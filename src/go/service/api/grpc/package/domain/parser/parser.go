package parser

import (
	. "github.com/tilau2328/x/src/go/package/x"
	"github.com/tilau2328/x/src/go/service/api/grpc/package/domain/lexer"
	"github.com/tilau2328/x/src/go/service/api/grpc/package/domain/lexer/scanner"
	"github.com/tilau2328/x/src/go/service/api/grpc/package/domain/model"
	"github.com/tilau2328/x/src/go/service/api/grpc/package/domain/model/meta"
	"strings"
	"unicode/utf8"
)

type Parser struct {
	lex                   *lexer.Lexer
	permissive            bool
	bodyIncludingComments bool
}

func Permissive(permissive bool) Opt[*Parser] {
	return func(p *Parser) { p.permissive = permissive }
}
func BodyIncludingComments(bodyIncludingComments bool) Opt[*Parser] {
	return func(p *Parser) { p.bodyIncludingComments = bodyIncludingComments }
}
func NewParser(lex *lexer.Lexer, opts ...Opt[*Parser]) *Parser {
	p := &Parser{lex: lex}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
func (p *Parser) IsEOF() bool {
	p.lex.Next()
	defer p.lex.UnNext()
	return p.lex.IsEOF()
}
func (p *Parser) ParseComments() []*model.Comment {
	var comments []*model.Comment
	for {
		comment, err := p.parseComment()
		if err != nil {
			return comments
		}
		comments = append(comments, comment)
	}
}
func (p *Parser) parseComment() (*model.Comment, error) {
	p.lex.NextComment()
	if p.lex.Token == scanner.TCOMMENT {
		return &model.Comment{
			Raw: p.lex.Text,
			Meta: meta.Meta{
				Pos:     p.lex.Pos.Position,
				LastPos: p.lex.Pos.AdvancedBulk(p.lex.Text).Position,
			},
		}, nil
	}
	defer p.lex.UnNext()
	return nil, p.unexpected("comment")
}

func (p *Parser) ParseSyntax() (*model.Syntax, error) {
	if p.lex.NextKeyword(); p.lex.Token != scanner.TSYNTAX {
		return nil, p.unexpected("syntax")
	}
	startPos := p.lex.Pos
	if p.lex.Next(); p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}
	if p.lex.Next(); p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}
	lq := p.lex.Text
	if p.lex.Next(); p.lex.Text != "proto3" && p.lex.Text != "proto2" {
		return nil, p.unexpected("proto3 or proto2")
	}
	version := p.lex.Text
	if p.lex.Next(); p.lex.Token != scanner.TQUOTE {
		return nil, p.unexpected("quote")
	}
	tq := p.lex.Text
	if p.lex.Next(); p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}
	return &model.Syntax{
		ProtobufVersion: version,
		VersionQuote:    lq + version + tq,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
func (p *Parser) ParseEnum() (*model.Enum, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TENUM {
		return nil, p.unexpected("enum")
	}
	startPos := p.lex.Pos
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("enumName")
	}
	enumName := p.lex.Text
	enumBody, inlineLeftCurly, lastPos, err := p.parseEnumBody()
	if err != nil {
		return nil, err
	}

	return &model.Enum{
		Name:       enumName,
		Body:       enumBody,
		CommentBLC: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}
func (p *Parser) parseEnumBody() (
	[]model.Visitee,
	*model.Comment,
	scanner.Position,
	error,
) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, scanner.Position{}, p.unexpected("{")
	}
	inlineLeftCurly := p.parseInlineComment()
	var stmts []model.Visitee
	for {
		comments := p.ParseComments()
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()
		var stmt interface {
			model.HasInlineCommentSetter
			model.Visitee
		}
		switch token {
		case scanner.TRIGHTCURLY:
			if p.bodyIncludingComments {
				for _, comment := range comments {
					stmts = append(stmts, model.Visitee(comment))
				}
			}
			p.lex.Next()
			lastPos := p.lex.Pos
			if p.permissive {
				// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
				p.lex.ConsumeToken(scanner.TSEMICOLON)
				if p.lex.Token == scanner.TSEMICOLON {
					lastPos = p.lex.Pos
				}
			}
			return stmts, inlineLeftCurly, lastPos, nil
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TRESERVED:
			reserved, err := p.ParseReserved()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			reserved.Comments = comments
			stmt = reserved
		default:
			enumField, enumFieldErr := p.parseEnumField()
			if enumFieldErr == nil {
				enumField.Comments = comments
				stmt = enumField
				break
			}
			p.lex.UnNext()
			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmt = &model.EmptyStatement{}
				break
			}
			return nil, nil, scanner.Position{}, &parseEnumBodyStatementErr{
				parseEnumFieldErr:      enumFieldErr,
				parseEmptyStatementErr: emptyErr,
			}
		}
		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}
func (p *Parser) parseEnumField() (*model.EnumField, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("ident")
	}
	startPos := p.lex.Pos
	ident := p.lex.Text
	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}
	var intLit string
	p.lex.ConsumeToken(scanner.TMINUS)
	if p.lex.Token == scanner.TMINUS {
		intLit = "-"
	}
	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		return nil, p.unexpected("intLit")
	}
	intLit += p.lex.Text
	enumValueOptions, err := p.parseEnumValueOptions()
	if err != nil {
		return nil, err
	}
	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}
	return &model.EnumField{
		Ident:   ident,
		Number:  intLit,
		Options: enumValueOptions,
		Meta:    meta.Meta{Pos: startPos.Position},
	}, nil
}
func (p *Parser) parseEnumValueOptions() ([]*model.EnumValueOption, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTSQUARE {
		p.lex.UnNext()
		return nil, nil
	}
	opt, err := p.parseEnumValueOption()
	if err != nil {
		return nil, p.unexpected("enumValueOption")
	}
	var opts []*model.EnumValueOption
	opts = append(opts, opt)
	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}
		opt, err = p.parseEnumValueOption()
		if err != nil {
			return nil, p.unexpected("enumValueOption")
		}
		opts = append(opts, opt)
	}
	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTSQUARE {
		return nil, p.unexpected("]")
	}
	return opts, nil
}

func (p *Parser) parseEnumValueOption() (*model.EnumValueOption, error) {
	optionName, err := p.parseOptionName()
	if err != nil {
		return nil, err
	}
	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}
	constant, err := p.parseOptionConstant()
	if err != nil {
		return nil, err
	}
	return &model.EnumValueOption{
		Name:     optionName,
		Constant: constant,
	}, nil
}
func (p *Parser) ParseExtensions() (*model.Extensions, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TEXTENSIONS {
		return nil, p.unexpected("extensions")
	}
	startPos := p.lex.Pos
	ranges, err := p.parseRanges()
	if err != nil {
		return nil, err
	}
	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}
	return &model.Extensions{
		Ranges: ranges,
		Meta:   meta.Meta{Pos: startPos.Position},
	}, nil
}

func (p *Parser) ParseExtend() (*model.Extend, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TEXTEND {
		return nil, p.unexpected("extend")
	}
	startPos := p.lex.Pos
	messageType, _, err := p.lex.ReadMessageType()
	if err != nil {
		return nil, err
	}
	extendBody, inlineLeftCurly, lastPos, err := p.parseExtendBody()
	if err != nil {
		return nil, err
	}
	return &model.Extend{
		Type:       messageType,
		Body:       extendBody,
		CommentBLC: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

func (p *Parser) parseExtendBody() ([]model.Visitee, *model.Comment, scanner.Position, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, scanner.Position{}, p.unexpected("{")
	}
	inlineLeftCurly := p.parseInlineComment()
	p.lex.Next()
	if p.lex.Token == scanner.TRIGHTCURLY {
		lastPos := p.lex.Pos
		if p.permissive {
			p.lex.ConsumeToken(scanner.TSEMICOLON)
			if p.lex.Token == scanner.TSEMICOLON {
				lastPos = p.lex.Pos
			}
		}
		return nil, nil, lastPos, nil
	}
	p.lex.UnNext()
	var stmts []model.Visitee
	for {
		comments := p.ParseComments()
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()
		var stmt interface {
			model.HasInlineCommentSetter
			model.Visitee
		}
		switch token {
		case scanner.TRIGHTCURLY:
			if p.bodyIncludingComments {
				for _, comment := range comments {
					stmts = append(stmts, model.Visitee(comment))
				}
			}
			p.lex.Next()
			lastPos := p.lex.Pos
			if p.permissive {
				p.lex.ConsumeToken(scanner.TSEMICOLON)
				if p.lex.Token == scanner.TSEMICOLON {
					lastPos = p.lex.Pos
				}
			}
			return stmts, inlineLeftCurly, lastPos, nil
		default:
			field, fieldErr := p.ParseField()
			if fieldErr == nil {
				field.Comments = comments
				stmt = field
				break
			}
			p.lex.UnNext()
			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmt = &model.EmptyStatement{}
				break
			}

			return nil, nil, scanner.Position{}, &parseExtendBodyStatementErr{
				parseFieldErr:          fieldErr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}
func (p *Parser) ParseField() (ret *model.Field, err error) {
	p.lex.NextKeyword()
	ret = &model.Field{Meta: meta.Meta{Pos: p.lex.Pos.Position}}
	switch p.lex.Token {
	case scanner.TREPEATED:
		ret.IsRepeated = true
	case scanner.TREQUIRED:
		ret.IsRequired = true
	case scanner.TOPTIONAL:
		ret.IsOptional = true
	default:
		p.lex.UnNext()
	}
	if ret.Type, _, err = p.parseType(); err != nil {
		return nil, p.unexpected("type")
	}
	p.lex.Next()
	ret.Name = p.lex.Text
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("fieldName")
	} else if p.lex.Next(); p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	} else if ret.Number, err = p.parseFieldNumber(); err != nil {
		return nil, p.unexpected("fieldNumber")
	} else if ret.Options, err = p.parseFieldOptionsOption(); err != nil {
		return nil, err
	} else if p.lex.Next(); p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}
	return
}
func (p *Parser) parseFieldOptionsOption() ([]*model.FieldOption, error) {
	if p.lex.Next(); p.lex.Token == scanner.TLEFTSQUARE {
		if fieldOptions, err := p.parseFieldOptions(); err != nil {
			return nil, err
		} else if p.lex.Next(); p.lex.Token != scanner.TRIGHTSQUARE {
			return nil, p.unexpected("]")
		} else {
			return fieldOptions, nil
		}
	}
	p.lex.UnNext()
	return nil, nil
}

func (p *Parser) parseFieldOptions() (opts []*model.FieldOption, err error) {
	opt, err := p.parseFieldOption()
	if err != nil {
		return nil, err
	}
	opts = append(opts, opt)
	for {
		if p.lex.Next(); p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		} else if opt, err = p.parseFieldOption(); err != nil {
			return nil, p.unexpected("fieldOption")
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

func (p *Parser) parseFieldOption() (*model.FieldOption, error) {
	if optionName, err := p.parseOptionName(); err != nil {
		return nil, err
	} else if p.lex.Next(); p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	} else if constant, err := p.parseOptionConstant(); err != nil {
		return nil, err
	} else {
		return &model.FieldOption{Name: optionName, Constant: constant}, nil
	}
}

var typeConstants = map[string]struct{}{
	"double":   {},
	"float":    {},
	"int32":    {},
	"int64":    {},
	"uint32":   {},
	"uint64":   {},
	"sint32":   {},
	"sint64":   {},
	"fixed32":  {},
	"fixed64":  {},
	"sfixed32": {},
	"sfixed64": {},
	"bool":     {},
	"string":   {},
	"bytes":    {},
}

func (p *Parser) parseType() (string, scanner.Position, error) {
	p.lex.Next()
	if _, ok := typeConstants[p.lex.Text]; ok {
		return p.lex.Text, p.lex.Pos, nil
	}
	p.lex.UnNext()
	messageOrEnumType, startPos, err := p.lex.ReadMessageType()
	if err != nil {
		return "", scanner.Position{}, err
	}
	return messageOrEnumType, startPos, nil
}

func (p *Parser) parseFieldNumber() (string, error) {
	if p.lex.NextNumberLit(); p.lex.Token != scanner.TINTLIT {
		return "", p.unexpected("intLit")
	}
	return p.lex.Text, nil
}
func (p *Parser) ParseGroupField() (*model.GroupField, error) {
	var isRepeated bool
	var isRequired bool
	var isOptional bool
	p.lex.NextKeyword()
	startPos := p.lex.Pos
	if p.lex.Token == scanner.TREPEATED {
		isRepeated = true
	} else if p.lex.Token == scanner.TREQUIRED {
		isRequired = true
	} else if p.lex.Token == scanner.TOPTIONAL {
		isOptional = true
	} else {
		p.lex.UnNext()
	}

	p.lex.NextKeyword()
	if p.lex.Token != scanner.TGROUP {
		return nil, p.unexpected("group")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("groupName")
	}
	if !isCapitalized(p.lex.Text) {
		return nil, p.unexpectedf("groupName %q must begin with capital letter.", p.lex.Text)
	}
	groupName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	fieldNumber, err := p.parseFieldNumber()
	if err != nil {
		return nil, p.unexpected("fieldNumber")
	}

	messageBody, inlineLeftCurly, lastPos, err := p.parseMessageBody()
	if err != nil {
		return nil, err
	}

	return &model.GroupField{
		IsRepeated: isRepeated,
		IsRequired: isRequired,
		IsOptional: isOptional,
		GroupName:  groupName,
		Number:     fieldNumber,
		Body:       messageBody,

		CommentBLC: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}
func (p *Parser) peekIsGroup() bool {
	p.lex.NextKeyword()
	switch p.lex.Token {
	case scanner.TREPEATED,
		scanner.TREQUIRED,
		scanner.TOPTIONAL:
		defer p.lex.UnNextTo(p.lex.RawText)
	default:
		p.lex.UnNext()
	}

	p.lex.NextKeyword()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TGROUP {
		return false
	}

	p.lex.Next()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TIDENT {
		return false
	}
	if !isCapitalized(p.lex.Text) {
		return false
	}

	p.lex.Next()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TEQUALS {
		return false
	}

	_, err := p.parseFieldNumber()
	defer p.lex.UnNextTo(p.lex.RawText)
	if err != nil {
		return false
	}

	p.lex.Next()
	defer p.lex.UnNextTo(p.lex.RawText)
	if p.lex.Token != scanner.TLEFTCURLY {
		return false
	}
	return true
}

func isCapitalized(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return isUpper(r)
}

func isUpper(r rune) bool {
	return 'From' <= r && r <= 'Z'
}

func (p *Parser) ParseImport() (*model.Import, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TIMPORT {
		return nil, p.unexpected(`"import"`)
	}
	startPos := p.lex.Pos

	var modifier model.ImportModifier
	p.lex.NextKeywordOrStrLit()
	switch p.lex.Token {
	case scanner.TPUBLIC:
		modifier = model.ImportModifierPublic
	case scanner.TWEAK:
		modifier = model.ImportModifierWeak
	case scanner.TSTRLIT:
		modifier = model.ImportModifierNone
		p.lex.UnNext()
	}

	p.lex.NextStrLit()
	if p.lex.Token != scanner.TSTRLIT {
		return nil, p.unexpected("strLit")
	}
	location := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &model.Import{
		Modifier: modifier,
		Location: location,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
func (p *Parser) MaybeScanInlineComment(hasSetter model.HasInlineCommentSetter) {
	inlineComment := p.parseInlineComment()
	if inlineComment == nil {
		return
	}
	hasSetter.SetComment(inlineComment)
}
func (p *Parser) parseInlineComment() *model.Comment {
	currentPos := p.lex.Pos
	comment, err := p.parseComment()
	if err != nil {
		return nil
	}
	if currentPos.Line != comment.Meta.Pos.Line {
		p.lex.UnNext()
		return nil
	}
	return comment
}
func (p *Parser) ParsePackage() (*model.Package, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TPACKAGE {
		return nil, p.unexpected("package")
	}
	startPos := p.lex.Pos

	ident, _, err := p.lex.ReadFullIdent()
	if err != nil {
		return nil, p.unexpected("fullIdent")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &model.Package{
		Name: ident,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
func (p *Parser) ParseMessage() (*model.Message, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TMESSAGE {
		return nil, p.unexpected("message")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("messageName")
	}
	messageName := p.lex.Text

	messageBody, inlineLeftCurly, lastPos, err := p.parseMessageBody()
	if err != nil {
		return nil, err
	}

	return &model.Message{
		MessageName: messageName,
		MessageBody: messageBody,
		CommentBLC:  inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

func (p *Parser) parseMessageBody() (
	[]model.Visitee,
	*model.Comment,
	scanner.Position,
	error,
) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, scanner.Position{}, p.unexpected("{")
	}
	inlineLeftCurly := p.parseInlineComment()
	p.lex.Next()
	if p.lex.Token == scanner.TRIGHTCURLY {
		lastPos := p.lex.Pos
		if p.permissive {
			// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
			p.lex.ConsumeToken(scanner.TSEMICOLON)
			if p.lex.Token == scanner.TSEMICOLON {
				lastPos = p.lex.Pos
			}
		}

		return nil, nil, lastPos, nil
	}
	p.lex.UnNext()
	var stmts []model.Visitee
	for {
		comments := p.ParseComments()
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()
		var stmt interface {
			model.HasInlineCommentSetter
			model.Visitee
		}
		switch token {
		case scanner.TRIGHTCURLY:
			if p.bodyIncludingComments {
				for _, comment := range comments {
					stmts = append(stmts, model.Visitee(comment))
				}
			}
			p.lex.Next()

			lastPos := p.lex.Pos
			if p.permissive {
				// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
				p.lex.ConsumeToken(scanner.TSEMICOLON)
				if p.lex.Token == scanner.TSEMICOLON {
					lastPos = p.lex.Pos
				}
			}
			return stmts, inlineLeftCurly, lastPos, nil
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			enum.Comments = comments
			stmt = enum
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			message.Comments = comments
			stmt = message
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TONEOF:
			oneof, err := p.ParseOneof()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			oneof.Comments = comments
			stmt = oneof
		case scanner.TMAP:
			mapField, err := p.ParseMapField()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			mapField.Comments = comments
			stmt = mapField
		case scanner.TEXTEND:
			extend, err := p.ParseExtend()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			extend.Comments = comments
			stmt = extend
		case scanner.TRESERVED:
			reserved, err := p.ParseReserved()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			reserved.Comments = comments
			stmt = reserved
		case scanner.TEXTENSIONS:
			extensions, err := p.ParseExtensions()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			extensions.Comments = comments
			stmt = extensions
		default:
			var ferr error
			isGroup := p.peekIsGroup()
			if isGroup {
				groupField, groupErr := p.ParseGroupField()
				if groupErr == nil {
					groupField.Comments = comments
					stmt = groupField
					break
				}
				ferr = groupErr
				p.lex.UnNext()
			} else {
				field, fieldErr := p.ParseField()
				if fieldErr == nil {
					field.Comments = comments
					stmt = field
					break
				}
				ferr = fieldErr
				p.lex.UnNext()
			}

			emptyErr := p.lex.ReadEmptyStatement()
			if emptyErr == nil {
				stmt = &model.EmptyStatement{}
				break
			}

			return nil, nil, scanner.Position{}, &parseMessageBodyStatementErr{
				parseFieldErr:          ferr,
				parseEmptyStatementErr: emptyErr,
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}

func (p *Parser) ParseMapField() (*model.MapField, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TMAP {
		return nil, p.unexpected("map")
	}
	startPos := p.lex.Pos
	p.lex.Next()
	if p.lex.Token != scanner.TLESS {
		return nil, p.unexpected("<")
	}
	keyType, err := p.parseKeyType()
	if err != nil {
		return nil, err
	}
	p.lex.Next()
	if p.lex.Token != scanner.TCOMMA {
		return nil, p.unexpected(",")
	}
	typeValue, _, err := p.parseType()
	if err != nil {
		return nil, p.unexpected("type")
	}
	p.lex.Next()
	if p.lex.Token != scanner.TGREATER {
		return nil, p.unexpected(">")
	}
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("mapName")
	}
	mapName := p.lex.Text
	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}
	fieldNumber, err := p.parseFieldNumber()
	if err != nil {
		return nil, p.unexpected("fieldNumber")
	}
	fieldOptions, err := p.parseFieldOptionsOption()
	if err != nil {
		return nil, err
	}
	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}
	return &model.MapField{
		KeyType:      keyType,
		Type:         typeValue,
		MapName:      mapName,
		FieldNumber:  fieldNumber,
		FieldOptions: fieldOptions,
		Meta:         meta.Meta{Pos: startPos.Position},
	}, nil
}

var keyTypeConstants = map[string]struct{}{
	"int32":    {},
	"int64":    {},
	"uint32":   {},
	"uint64":   {},
	"sint32":   {},
	"sint64":   {},
	"fixed32":  {},
	"fixed64":  {},
	"sfixed32": {},
	"sfixed64": {},
	"bool":     {},
	"string":   {},
}

func (p *Parser) parseKeyType() (string, error) {
	p.lex.Next()
	if _, ok := keyTypeConstants[p.lex.Text]; ok {
		return p.lex.Text, nil
	}
	return "", p.unexpected("keyType constant")
}

func (p *Parser) ParseService() (*model.Service, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TSERVICE {
		return nil, p.unexpected("service")
	}
	startPos := p.lex.Pos
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("serviceName")
	}
	serviceName := p.lex.Text
	serviceBody, inlineLeftCurly, lastPos, err := p.parseServiceBody()
	if err != nil {
		return nil, err
	}
	return &model.Service{
		Name:       serviceName,
		Body:       serviceBody,
		CommentBLC: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}

func (p *Parser) parseServiceBody() (
	[]model.Visitee,
	*model.Comment,
	scanner.Position,
	error,
) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, scanner.Position{}, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var stmts []model.Visitee
	for {
		comments := p.ParseComments()

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		var stmt interface {
			model.HasInlineCommentSetter
			model.Visitee
		}

		switch token {
		case scanner.TRIGHTCURLY:
			if p.bodyIncludingComments {
				for _, comment := range comments {
					stmts = append(stmts, model.Visitee(comment))
				}
			}
			p.lex.Next()

			lastPos := p.lex.Pos
			if p.permissive {
				// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
				p.lex.ConsumeToken(scanner.TSEMICOLON)
				if p.lex.Token == scanner.TSEMICOLON {
					lastPos = p.lex.Pos
				}
			}
			return stmts, inlineLeftCurly, lastPos, nil
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TRPC:
			rpc, err := p.parseRPC()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
			rpc.Comments = comments
			stmt = rpc
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, nil, scanner.Position{}, err
			}
		}

		p.MaybeScanInlineComment(stmt)
		stmts = append(stmts, stmt)
	}
}
func (p *Parser) parseRPC() (*model.RPC, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TRPC {
		return nil, p.unexpected("rpc")
	}
	startPos := p.lex.Pos

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("serviceName")
	}
	rpcName := p.lex.Text

	rpcRequest, err := p.parseRPCRequest()
	if err != nil {
		return nil, err
	}

	p.lex.NextKeyword()
	if p.lex.Token != scanner.TRETURNS {
		return nil, p.unexpected("returns")
	}

	rpcResponse, err := p.parseRPCResponse()
	if err != nil {
		return nil, err
	}

	var opts []*model.Option
	var inlineLeftCurly *model.Comment
	p.lex.Next()
	lastPos := p.lex.Pos
	switch p.lex.Token {
	case scanner.TLEFTCURLY:
		p.lex.UnNext()
		opts, inlineLeftCurly, err = p.parseRPCOptions()
		if err != nil {
			return nil, err
		}
		lastPos = p.lex.Pos
		if p.permissive {
			// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
			p.lex.ConsumeToken(scanner.TSEMICOLON)
			if p.lex.Token == scanner.TSEMICOLON {
				lastPos = p.lex.Pos
			}
		}
	case scanner.TSEMICOLON:
		break
	default:
		return nil, p.unexpected("{ or ;")
	}

	return &model.RPC{
		Name:       rpcName,
		Request:    rpcRequest,
		Response:   rpcResponse,
		Options:    opts,
		CommentBLC: inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}
func (p *Parser) parseRPCRequest() (*model.Request, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTPAREN {
		return nil, p.unexpected("(")
	}
	startPos := p.lex.Pos

	p.lex.NextKeyword()
	isStream := true
	if p.lex.Token != scanner.TSTREAM {
		isStream = false
		p.lex.UnNext()
	}

	messageType, _, err := p.lex.ReadMessageType()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTPAREN {
		return nil, p.unexpected(")")
	}

	return &model.Request{
		IsStream:    isStream,
		MessageType: messageType,
		Meta:        meta.Meta{Pos: startPos.Position},
	}, nil
}
func (p *Parser) parseRPCResponse() (*model.Response, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTPAREN {
		return nil, p.unexpected("(")
	}
	startPos := p.lex.Pos

	p.lex.NextKeyword()
	isStream := true
	if p.lex.Token != scanner.TSTREAM {
		isStream = false
		p.lex.UnNext()
	}

	messageType, _, err := p.lex.ReadMessageType()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TRIGHTPAREN {
		return nil, p.unexpected(")")
	}

	return &model.Response{
		IsStream: isStream,
		Type:     messageType,
		Meta:     meta.Meta{Pos: startPos.Position},
	}, nil
}
func (p *Parser) parseRPCOptions() ([]*model.Option, *model.Comment, error) {
	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, nil, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var options []*model.Option
	for {
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()

		switch token {
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, nil, err
			}
			options = append(options, option)
		case scanner.TRIGHTCURLY:
			// This spec is not documented, but allowed in general.
			break
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, nil, err
			}
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			return options, inlineLeftCurly, nil
		}
		p.lex.UnNext()
	}
}
func (p *Parser) ParseOption() (*model.Option, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TOPTION {
		return nil, p.unexpected("option")
	}
	startPos := p.lex.Pos

	optionName, err := p.parseOptionName()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	constant, err := p.parseOptionConstant()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &model.Option{
		Name:     optionName,
		Constant: constant,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: p.lex.Pos.Position,
		},
	}, nil
}
func (p *Parser) parseCloudEndpointsOptionConstant() (string, error) {
	var ret string

	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return "", p.unexpected("{")
	}
	ret += p.lex.Text

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TIDENT {
			return "", p.unexpected("ident")
		}
		ret += p.lex.Text

		needSemi := false
		p.lex.Next()
		switch p.lex.Token {
		case scanner.TLEFTCURLY:
			if !p.permissive {
				return "", p.unexpected(":")
			}
			p.lex.UnNext()
		case scanner.TCOLON:
			ret += p.lex.Text
			if p.lex.Peek() == scanner.TLEFTCURLY && p.permissive {
				needSemi = true
			}
		default:
			if p.permissive {
				return "", p.unexpected("{ or :")
			}
			return "", p.unexpected(":")
		}

		constant, err := p.parseOptionConstant()
		if err != nil {
			return "", err
		}
		ret += constant

		p.lex.Next()
		if p.lex.Token == scanner.TSEMICOLON && needSemi && p.permissive {
			ret += p.lex.Text
			p.lex.Next()
		}

		switch {
		case p.lex.Token == scanner.TCOMMA, p.lex.Token == scanner.TSEMICOLON:
			ret += p.lex.Text
			if p.lex.Peek() == scanner.TRIGHTCURLY && p.permissive {
				p.lex.Next()
				ret += p.lex.Text
				return ret, nil
			}
		case p.lex.Token == scanner.TRIGHTCURLY:
			ret += p.lex.Text
			return ret, nil
		default:
			ret += "\n"
			p.lex.UnNext()
		}
	}
}
func (p *Parser) parseOptionName() (string, error) {
	var optionName string
	p.lex.Next()
	switch p.lex.Token {
	case scanner.TIDENT:
		optionName = p.lex.Text
	case scanner.TLEFTPAREN:
		optionName = p.lex.Text

		// protoc accepts "(." fullIndent ")". See #63
		if p.permissive {
			p.lex.Next()
			if p.lex.Token == scanner.TDOT {
				optionName += "."
			} else {
				p.lex.UnNext()
			}
		}

		fullIdent, _, err := p.lex.ReadFullIdent()
		if err != nil {
			return "", err
		}
		optionName += fullIdent

		p.lex.Next()
		if p.lex.Token != scanner.TRIGHTPAREN {
			return "", p.unexpected(")")
		}
		optionName += p.lex.Text
	default:
		return "", p.unexpected("ident or left paren")
	}

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TDOT {
			p.lex.UnNext()
			break
		}
		optionName += p.lex.Text

		p.lex.Next()
		if p.lex.Token != scanner.TIDENT {
			return "", p.unexpected("ident")
		}
		optionName += p.lex.Text
	}
	return optionName, nil
}
func (p *Parser) parseOptionConstant() (constant string, err error) {
	switch p.lex.Peek() {
	case scanner.TLEFTCURLY:
		if !p.permissive {
			return "", p.unexpected("constant or permissive mode")
		}

		// parses empty fields within an option
		if p.lex.PeekN(2) == scanner.TRIGHTCURLY {
			p.lex.NextN(2)
			return "{}", nil
		}

		constant, err = p.parseCloudEndpointsOptionConstant()
		if err != nil {
			return "", err
		}

	case scanner.TLEFTSQUARE:
		if !p.permissive {
			return "", p.unexpected("constant or permissive mode")
		}
		p.lex.Next()

		// parses empty fields within an option
		if p.lex.Peek() == scanner.TRIGHTSQUARE {
			p.lex.Next()
			return "[]", nil
		}

		constant, err = p.parseOptionConstants()
		if err != nil {
			return "", err
		}
		p.lex.Next()
		constant = "[" + constant + "]"

	default:
		constant, _, err = p.lex.ReadConstant(p.permissive)
		if err != nil {
			return "", err
		}
	}
	return constant, nil
}
func (p *Parser) parseOptionConstants() (constant string, err error) {
	opt, err := p.parseOptionConstant()
	if err != nil {
		return "", err
	}
	var opts []string
	opts = append(opts, opt)
	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}
		opt, err = p.parseOptionConstant()
		if err != nil {
			return "", p.unexpected("optionConstant")
		}
		opts = append(opts, opt)
	}
	return strings.Join(opts, ","), nil
}

func (p *Parser) ParseOneof() (*model.OneOf, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TONEOF {
		return nil, p.unexpected("oneof")
	}
	startPos := p.lex.Pos
	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("oneofName")
	}
	oneofName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TLEFTCURLY {
		return nil, p.unexpected("{")
	}

	inlineLeftCurly := p.parseInlineComment()

	var oneofFields []*model.OneOfField
	var options []*model.Option
	for {
		comments := p.ParseComments()

		err := p.lex.ReadEmptyStatement()
		if err == nil {
			continue
		}

		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()
		if token == scanner.TOPTION {
			// See https://github.com/yoheimuta/go-protoparser/issues/57
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			option.Comments = comments
			p.MaybeScanInlineComment(option)
			options = append(options, option)
		} else {
			oneofField, err := p.parseOneofField()
			if err != nil {
				return nil, err
			}
			oneofField.Comments = comments
			p.MaybeScanInlineComment(oneofField)
			oneofFields = append(oneofFields, oneofField)
		}

		p.lex.Next()
		if p.lex.Token == scanner.TRIGHTCURLY {
			break
		} else {
			p.lex.UnNext()
		}
	}

	lastPos := p.lex.Pos
	if p.permissive {
		// accept a block followed by semicolon. See https://github.com/yoheimuta/go-protoparser/v4/issues/30.
		p.lex.ConsumeToken(scanner.TSEMICOLON)
		if p.lex.Token == scanner.TSEMICOLON {
			lastPos = p.lex.Pos
		}
	}

	return &model.OneOf{
		OneofFields: oneofFields,
		OneofName:   oneofName,
		Options:     options,
		CommentBLC:  inlineLeftCurly,
		Meta: meta.Meta{
			Pos:     startPos.Position,
			LastPos: lastPos.Position,
		},
	}, nil
}
func (p *Parser) parseOneofField() (*model.OneOfField, error) {
	typeValue, startPos, err := p.parseType()
	if err != nil {
		return nil, p.unexpected("type")
	}

	p.lex.Next()
	if p.lex.Token != scanner.TIDENT {
		return nil, p.unexpected("fieldName")
	}
	fieldName := p.lex.Text

	p.lex.Next()
	if p.lex.Token != scanner.TEQUALS {
		return nil, p.unexpected("=")
	}

	fieldNumber, err := p.parseFieldNumber()
	if err != nil {
		return nil, p.unexpected("fieldNumber")
	}

	fieldOptions, err := p.parseFieldOptionsOption()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}

	return &model.OneOfField{
		Type:    typeValue,
		Name:    fieldName,
		Number:  fieldNumber,
		Options: fieldOptions,
		Meta:    meta.Meta{Pos: startPos.Position},
	}, nil
}

func (p *Parser) ParseProto() (*model.Proto, error) {
	syntaxComments := p.ParseComments()
	syntax, err := p.ParseSyntax()
	if err != nil {
		return nil, err
	}
	syntax.Comments = syntaxComments
	p.MaybeScanInlineComment(syntax)

	protoBody, err := p.parseProtoBody()
	if err != nil {
		return nil, err
	}

	return &model.Proto{
		Syntax:    syntax,
		ProtoBody: protoBody,
		Meta: &model.ProtoMeta{
			Filename: p.lex.Pos.Filename,
		},
	}, nil
}
func (p *Parser) parseProtoBody() ([]model.Visitee, error) {
	var protoBody []model.Visitee
	for {
		comments := p.ParseComments()
		if p.IsEOF() {
			if p.bodyIncludingComments {
				for _, comment := range comments {
					protoBody = append(protoBody, model.Visitee(comment))
				}
			}
			return protoBody, nil
		}
		p.lex.NextKeyword()
		token := p.lex.Token
		p.lex.UnNext()
		var stmt interface {
			model.HasInlineCommentSetter
			model.Visitee
		}
		switch token {
		case scanner.TIMPORT:
			importValue, err := p.ParseImport()
			if err != nil {
				return nil, err
			}
			importValue.Comments = comments
			stmt = importValue
		case scanner.TPACKAGE:
			packageValue, err := p.ParsePackage()
			if err != nil {
				return nil, err
			}
			packageValue.Comments = comments
			stmt = packageValue
		case scanner.TOPTION:
			option, err := p.ParseOption()
			if err != nil {
				return nil, err
			}
			option.Comments = comments
			stmt = option
		case scanner.TMESSAGE:
			message, err := p.ParseMessage()
			if err != nil {
				return nil, err
			}
			message.Comments = comments
			stmt = message
		case scanner.TENUM:
			enum, err := p.ParseEnum()
			if err != nil {
				return nil, err
			}
			enum.Comments = comments
			stmt = enum
		case scanner.TSERVICE:
			service, err := p.ParseService()
			if err != nil {
				return nil, err
			}
			service.Comments = comments
			stmt = service
		case scanner.TEXTEND:
			extend, err := p.ParseExtend()
			if err != nil {
				return nil, err
			}
			extend.Comments = comments
			stmt = extend
		default:
			err := p.lex.ReadEmptyStatement()
			if err != nil {
				return nil, err
			}
			protoBody = append(protoBody, &model.EmptyStatement{})
		}
		p.MaybeScanInlineComment(stmt)
		protoBody = append(protoBody, stmt)
	}
}
func (p *Parser) ParseReserved() (*model.Reserved, error) {
	p.lex.NextKeyword()
	if p.lex.Token != scanner.TRESERVED {
		return nil, p.unexpected("reserved")
	}
	startPos := p.lex.Pos

	parse := func() ([]*model.Range, []string, error) {
		ranges, err := p.parseRanges()
		if err == nil {
			return ranges, nil, nil
		}
		fieldNames, ferr := p.parseFieldNames()
		if ferr == nil {
			return nil, fieldNames, nil
		}
		return nil, nil, &parseReservedErr{
			parseRangesErr:     err,
			parseFieldNamesErr: ferr,
		}
	}

	ranges, fieldNames, err := parse()
	if err != nil {
		return nil, err
	}

	p.lex.Next()
	if p.lex.Token != scanner.TSEMICOLON {
		return nil, p.unexpected(";")
	}
	return &model.Reserved{
		Ranges:     ranges,
		FieldNames: fieldNames,
		Meta:       meta.Meta{Pos: startPos.Position},
	}, nil
}
func (p *Parser) parseRanges() ([]*model.Range, error) {
	var ranges []*model.Range
	rangeValue, err := p.parseRange()
	if err != nil {
		return nil, err
	}
	ranges = append(ranges, rangeValue)

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}

		rangeValue, err := p.parseRange()
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, rangeValue)
	}
	return ranges, nil
}
func (p *Parser) parseRange() (*model.Range, error) {
	p.lex.NextNumberLit()
	if p.lex.Token != scanner.TINTLIT {
		p.lex.UnNext()
		return nil, p.unexpected("intLit")
	}
	begin := p.lex.Text

	p.lex.Next()
	if p.lex.Text != "to" {
		p.lex.UnNext()
		return &model.Range{
			Begin: begin,
		}, nil
	}

	p.lex.NextNumberLit()
	switch {
	case p.lex.Token == scanner.TINTLIT,
		p.lex.Text == "max":
		return &model.Range{
			Begin: begin,
			End:   p.lex.Text,
		}, nil
	default:
		break
	}
	return nil, p.unexpected(`"intLit | "max"`)
}
func (p *Parser) parseFieldNames() ([]string, error) {
	var fieldNames []string

	fieldName, err := p.parseQuotedFieldName()
	if err != nil {
		return nil, err
	}
	fieldNames = append(fieldNames, fieldName)

	for {
		p.lex.Next()
		if p.lex.Token != scanner.TCOMMA {
			p.lex.UnNext()
			break
		}

		fieldName, err = p.parseQuotedFieldName()
		if err != nil {
			return nil, err
		}
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames, nil
}
func (p *Parser) parseQuotedFieldName() (string, error) {
	p.lex.NextStrLit()
	if p.lex.Token != scanner.TSTRLIT {
		p.lex.UnNext()
		return "", p.unexpected("quotedFieldName")
	}
	return p.lex.Text, nil
}
