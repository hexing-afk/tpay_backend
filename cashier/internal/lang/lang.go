package lang

const (
	// 规范化语言类型
	LangEnUS = "en_US"
	LangKmKH = "km_KH"
	LangZhCN = "zh_CN"
	LangThTH = "th_TH"

	// 系统默认语言为英语
	DefaultLang = LangEnUS
)

// 所有的语言列表
var LangList = []string{LangZhCN, LangEnUS, LangKmKH, LangThTH}

// 初始化各种语言
func init() {
	initLangEnUs() // 初始化英语
	initLangKmKH() // 初始化高棉语
	initLangThTH() // 初始化泰语
	initLangZhCN() // 初始化中文
}

func Lang(language, key string) string {
	switch language {
	case LangZhCN:
		if val, exist := LanguageZhCN[key]; exist {
			return val
		}
	case LangEnUS:
		if val, exist := LanguageEnUs[key]; exist {
			return val
		}
	case LangKmKH:
		if val, exist := LanguageKmKH[key]; exist {
			return val
		}
	case LangThTH:
		if val, exist := LanguageThTH[key]; exist {
			return val
		}
	}

	return key
}
