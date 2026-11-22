package xhtml

// Names of common HTML elements, to be used with filters such as
// FilterByElement.
const (
	TagHTML     = "html"
	TagTitle    = "title"
	TagHead     = "head"
	TagBody     = "body"
	TagHeader   = "header"
	TagFooter   = "footer"
	TagMain     = "main"
	TagSection  = "section"
	TagDiv      = "div"
	TagSpan     = "span"
	TagP        = "p"
	TagA        = "a"
	TagBr       = "br"
	TagHr       = "hr"
	TagUl       = "ul"
	TagLi       = "li"
	TagH1       = "h1"
	TagH2       = "h2"
	TagH3       = "h3"
	TagH4       = "h4"
	TagH5       = "h5"
	TagH6       = "h6"
	TagStrong   = "strong"
	TagEm       = "em"
	TagSmall    = "small"
	TagMark     = "mark"
	TagDel      = "del"
	TagIns      = "ins"
	TagImg      = "img"
	TagFigure   = "figure"
	TagTable    = "table"
	TagTr       = "tr"
	TagTd       = "td"
	TagTh       = "th"
	TagForm     = "form"
	TagInput    = "input"
	TagTextarea = "textarea"
	TagButton   = "button"
	TagSelect   = "select"
	TagOption   = "option"
)

// Names of common HTML attributes, to be used with filters such as
// FilterByAttribute or with GetAttributeValue.
const (
	AttrClass = "class"
	AttrHref  = "href"
	AttrID    = "id"
	AttrRef   = "ref"
	AttrSrc   = "src"
)
