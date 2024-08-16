package kana_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wantedly/kana-go"
)

func TestCanonicalConvert(t *testing.T) {
	var testcases = []struct {
		name    string
		input   string
		options kana.ConvertOptions
		expect  string
	}{
		{
			name:    "Base case ASCII Printable",
			input:   " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
			options: 0,
			expect:  " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		{
			name:    "Base case Latin-1 Printable",
			input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
			options: 0,
			expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		},
		{
			name:    "Base case Latin-1 Formatting",
			input:   "\u00A0\u00AD",
			options: 0,
			expect:  "\u00A0\u00AD",
		},
		{
			name:    "Base case General Punctuation Printable",
			input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
			options: 0,
			expect:  "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		},
		{
			name:    "Base case Some of Mathematical Operators",
			input:   "−∥",
			options: 0,
			expect:  "−∥",
		},
		{
			name:    "Base case CJK Symbols and Punctuation",
			input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
			options: 0,
			expect:  "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		},
		{
			name:    "Base case Hiragana",
			input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
			options: 0,
			expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		},
		{
			name:    "Base case Katakana",
			input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
			options: 0,
			expect:  "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		},
		{
			name:    "Base case Katakana Phonetic Extensions",
			input:   "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
			options: 0,
			expect:  "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
		},
		{
			name:    "Base case Fullwidth forms",
			input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
			options: 0,
			expect:  "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		},
		{
			name:    "Base case Halfwidth forms",
			input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
			options: 0,
			expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
		},
		{
			name:    "Base case Small Kana Extension",
			input:   "𛄲𛅐𛅑𛅒𛅕𛅤𛅥𛅦𛅧",
			options: 0,
			expect:  "𛄲𛅐𛅑𛅒𛅕𛅤𛅥𛅦𛅧",
		},
		{
			name:    "With KatakanaToHiragana Hiragana",
			input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
			options: kana.KatakanaToHiragana,
			expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		},
		{
			name:    "With KatakanaToHiragana Katakana",
			input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
			options: kana.KatakanaToHiragana,
			expect:  "゠ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖわ\u3099ゐ\u3099ゑ\u3099を\u3099・ーゝゞヿ",
		},
		{
			name:    "With KatakanaToHiragana Katakana Phonetic Extensions",
			input:   "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
			options: kana.KatakanaToHiragana,
			expect:  "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
		},
		{
			name:    "With KatakanaToHiragana Halfwidth Forms",
			input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
			options: kana.KatakanaToHiragana,
			expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
		},
		{
			name:    "With KatakanaToHiragana Small Kana Extension",
			input:   "𛄲𛅐𛅑𛅒𛅕𛅤𛅥𛅦𛅧",
			options: kana.KatakanaToHiragana,
			expect:  "𛄲𛅐𛅑𛅒𛄲𛅐𛅑𛅒𛅧",
		},
		{
			name:    "With HiraganaToKatakana Hiragana",
			input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
			options: kana.HiraganaToKatakana,
			expect:  "ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶ\u3099\u309A゛゜ヽヾゟ",
		},
		{
			name:    "With HiraganaToKatakana Katakana",
			input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
			options: kana.HiraganaToKatakana,
			expect:  "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		},
		{
			name:    "With HiraganaToKatakana Katakana Phonetic Extensions",
			input:   "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
			options: kana.HiraganaToKatakana,
			expect:  "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
		},
		{
			name:    "With HiraganaToKatakana Halfwidth Katakana",
			input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
			options: kana.HiraganaToKatakana,
			expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
		},
		{
			name:    "With HiraganaToKatakana Small Kana Extension",
			input:   "𛄲𛅐𛅑𛅒𛅕𛅤𛅥𛅦𛅧",
			options: kana.HiraganaToKatakana,
			expect:  "𛅕𛅤𛅥𛅦𛅕𛅤𛅥𛅦𛅧",
		},
		{
			name:    "With (KatakanaToHiragana | HiraganaToKatakana) Hiragana",
			input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
			options: kana.KatakanaToHiragana | kana.HiraganaToKatakana,
			expect:  "ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶ\u3099\u309A゛゜ヽヾゟ",
		},
		{
			name:    "With (KatakanaToHiragana | HiraganaToKatakana) Katakana",
			input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
			options: kana.KatakanaToHiragana | kana.HiraganaToKatakana,
			expect:  "゠ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖわ\u3099ゐ\u3099ゑ\u3099を\u3099・ーゝゞヿ",
		},
		{
			name:    "With (KatakanaToHiragana | HiraganaToKatakana) Katakana Phonetic Extensions",
			input:   "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
			options: kana.KatakanaToHiragana | kana.HiraganaToKatakana,
			expect:  "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
		},
		{
			name:    "With (KatakanaToHiragana | HiraganaToKatakana) Halfwidth forms",
			input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
			options: kana.KatakanaToHiragana | kana.HiraganaToKatakana,
			expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
		},
		{
			name:    "With (KatakanaToHiragana | HiraganaToKatakana) Small Kana Extension",
			input:   "𛄲𛅐𛅑𛅒𛅕𛅤𛅥𛅦𛅧",
			options: kana.KatakanaToHiragana | kana.HiraganaToKatakana,
			expect:  "𛅕𛅤𛅥𛅦𛄲𛅐𛅑𛅒𛅧",
		},
		{
			name:    "With HalfKanaToFull Halfwidth forms",
			input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
			options: kana.HalfwidthToWide,
			expect:  "。「」、・ヲァィゥェォャュョッーアイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワン\u3099\u309A\u3164ㄱㄲㄳㄴㄵㄶㄷㄸㄹㄺㄻㄼㄽㄾㄿㅀㅁㅂㅃㅄㅅㅆㅇㅈㅉㅊㅋㅌㅍㅎㅏㅐㅑㅒㅓㅔㅕㅖㅗㅘㅙㅚㅛㅜㅝㅞㅟㅠㅡㅢㅣ│←↑→↓■○",
		},
		{
			name:    "With HalfKanaToFull Voiced Composites",
			input:   "ｦ\uFF9Eｧ\uFF9Eｨ\uFF9Eｩ\uFF9Eｪ\uFF9Eｫ\uFF9Eｬ\uFF9Eｭ\uFF9Eｮ\uFF9Eｯ\uFF9Eｰ\uFF9Eｱ\uFF9Eｲ\uFF9Eｳ\uFF9Eｴ\uFF9Eｵ\uFF9Eｶ\uFF9Eｷ\uFF9Eｸ\uFF9Eｹ\uFF9Eｺ\uFF9Eｻ\uFF9Eｼ\uFF9Eｽ\uFF9Eｾ\uFF9Eｿ\uFF9Eﾀ\uFF9Eﾁ\uFF9Eﾂ\uFF9Eﾃ\uFF9Eﾄ\uFF9Eﾅ\uFF9Eﾆ\uFF9Eﾇ\uFF9Eﾈ\uFF9Eﾉ\uFF9Eﾊ\uFF9Eﾋ\uFF9Eﾌ\uFF9Eﾍ\uFF9Eﾎ\uFF9Eﾏ\uFF9Eﾐ\uFF9Eﾑ\uFF9Eﾒ\uFF9Eﾓ\uFF9Eﾔ\uFF9Eﾕ\uFF9Eﾖ\uFF9Eﾗ\uFF9Eﾘ\uFF9Eﾙ\uFF9Eﾚ\uFF9Eﾛ\uFF9Eﾜ\uFF9Eﾝ\uFF9E",
			options: kana.HalfwidthToWide,
			expect:  "ヺァ\u3099ィ\u3099ゥ\u3099ェ\u3099ォ\u3099ャ\u3099ュ\u3099ョ\u3099ッ\u3099ー\u3099ア\u3099イ\u3099ヴエ\u3099オ\u3099ガギグゲゴザジズゼゾダヂヅデドナ\u3099ニ\u3099ヌ\u3099ネ\u3099ノ\u3099バビブベボマ\u3099ミ\u3099ム\u3099メ\u3099モ\u3099ヤ\u3099ユ\u3099ヨ\u3099ラ\u3099リ\u3099ル\u3099レ\u3099ロ\u3099ヷン\u3099",
		},
		{
			name:    "With HalfKanaToFull Semi-Voiced Composites",
			input:   "ｦ\uFF9Fｧ\uFF9Fｨ\uFF9Fｩ\uFF9Fｪ\uFF9Fｫ\uFF9Fｬ\uFF9Fｭ\uFF9Fｮ\uFF9Fｯ\uFF9Fｰ\uFF9Fｱ\uFF9Fｲ\uFF9Fｳ\uFF9Fｴ\uFF9Fｵ\uFF9Fｶ\uFF9Fｷ\uFF9Fｸ\uFF9Fｹ\uFF9Fｺ\uFF9Fｻ\uFF9Fｼ\uFF9Fｽ\uFF9Fｾ\uFF9Fｿ\uFF9Fﾀ\uFF9Fﾁ\uFF9Fﾂ\uFF9Fﾃ\uFF9Fﾄ\uFF9Fﾅ\uFF9Fﾆ\uFF9Fﾇ\uFF9Fﾈ\uFF9Fﾉ\uFF9Fﾊ\uFF9Fﾋ\uFF9Fﾌ\uFF9Fﾍ\uFF9Fﾎ\uFF9Fﾏ\uFF9Fﾐ\uFF9Fﾑ\uFF9Fﾒ\uFF9Fﾓ\uFF9Fﾔ\uFF9Fﾕ\uFF9Fﾖ\uFF9Fﾗ\uFF9Fﾘ\uFF9Fﾙ\uFF9Fﾚ\uFF9Fﾛ\uFF9Fﾜ\uFF9Fﾝ\uFF9F",
			options: kana.HalfwidthToWide,
			expect:  "ヲ\u309Aァ\u309Aィ\u309Aゥ\u309Aェ\u309Aォ\u309Aャ\u309Aュ\u309Aョ\u309Aッ\u309Aー\u309Aア\u309Aイ\u309Aウ\u309Aエ\u309Aオ\u309Aカ\u309Aキ\u309Aク\u309Aケ\u309Aコ\u309Aサ\u309Aシ\u309Aス\u309Aセ\u309Aソ\u309Aタ\u309Aチ\u309Aツ\u309Aテ\u309Aト\u309Aナ\u309Aニ\u309Aヌ\u309Aネ\u309Aノ\u309Aパピプペポマ\u309Aミ\u309Aム\u309Aメ\u309Aモ\u309Aヤ\u309Aユ\u309Aヨ\u309Aラ\u309Aリ\u309Aル\u309Aレ\u309Aロ\u309Aワ\u309Aン\u309A",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Hiragana",
			input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Katakana",
			input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "゠ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖわ\u3099ゐ\u3099ゑ\u3099を\u3099・ーゝゞヿ",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Katakana Phonetic Extensions",
			input:   "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Halfwidth forms",
			input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "。「」、・をぁぃぅぇぉゃゅょっーあいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわん\u3099\u309A\u3164ㄱㄲㄳㄴㄵㄶㄷㄸㄹㄺㄻㄼㄽㄾㄿㅀㅁㅂㅃㅄㅅㅆㅇㅈㅉㅊㅋㅌㅍㅎㅏㅐㅑㅒㅓㅔㅕㅖㅗㅘㅙㅚㅛㅜㅝㅞㅟㅠㅡㅢㅣ│←↑→↓■○",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Voiced Composites",
			input:   "ｦ\uFF9Eｧ\uFF9Eｨ\uFF9Eｩ\uFF9Eｪ\uFF9Eｫ\uFF9Eｬ\uFF9Eｭ\uFF9Eｮ\uFF9Eｯ\uFF9Eｰ\uFF9Eｱ\uFF9Eｲ\uFF9Eｳ\uFF9Eｴ\uFF9Eｵ\uFF9Eｶ\uFF9Eｷ\uFF9Eｸ\uFF9Eｹ\uFF9Eｺ\uFF9Eｻ\uFF9Eｼ\uFF9Eｽ\uFF9Eｾ\uFF9Eｿ\uFF9Eﾀ\uFF9Eﾁ\uFF9Eﾂ\uFF9Eﾃ\uFF9Eﾄ\uFF9Eﾅ\uFF9Eﾆ\uFF9Eﾇ\uFF9Eﾈ\uFF9Eﾉ\uFF9Eﾊ\uFF9Eﾋ\uFF9Eﾌ\uFF9Eﾍ\uFF9Eﾎ\uFF9Eﾏ\uFF9Eﾐ\uFF9Eﾑ\uFF9Eﾒ\uFF9Eﾓ\uFF9Eﾔ\uFF9Eﾕ\uFF9Eﾖ\uFF9Eﾗ\uFF9Eﾘ\uFF9Eﾙ\uFF9Eﾚ\uFF9Eﾛ\uFF9Eﾜ\uFF9Eﾝ\uFF9E",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "を\u3099ぁ\u3099ぃ\u3099ぅ\u3099ぇ\u3099ぉ\u3099ゃ\u3099ゅ\u3099ょ\u3099っ\u3099ー\u3099あ\u3099い\u3099ゔえ\u3099お\u3099がぎぐげござじずぜぞだぢづでどな\u3099に\u3099ぬ\u3099ね\u3099の\u3099ばびぶべぼま\u3099み\u3099む\u3099め\u3099も\u3099や\u3099ゆ\u3099よ\u3099ら\u3099り\u3099る\u3099れ\u3099ろ\u3099わ\u3099ん\u3099",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Semi-Voiced Composites",
			input:   "ｦ\uFF9Fｧ\uFF9Fｨ\uFF9Fｩ\uFF9Fｪ\uFF9Fｫ\uFF9Fｬ\uFF9Fｭ\uFF9Fｮ\uFF9Fｯ\uFF9Fｰ\uFF9Fｱ\uFF9Fｲ\uFF9Fｳ\uFF9Fｴ\uFF9Fｵ\uFF9Fｶ\uFF9Fｷ\uFF9Fｸ\uFF9Fｹ\uFF9Fｺ\uFF9Fｻ\uFF9Fｼ\uFF9Fｽ\uFF9Fｾ\uFF9Fｿ\uFF9Fﾀ\uFF9Fﾁ\uFF9Fﾂ\uFF9Fﾃ\uFF9Fﾄ\uFF9Fﾅ\uFF9Fﾆ\uFF9Fﾇ\uFF9Fﾈ\uFF9Fﾉ\uFF9Fﾊ\uFF9Fﾋ\uFF9Fﾌ\uFF9Fﾍ\uFF9Fﾎ\uFF9Fﾏ\uFF9Fﾐ\uFF9Fﾑ\uFF9Fﾒ\uFF9Fﾓ\uFF9Fﾔ\uFF9Fﾕ\uFF9Fﾖ\uFF9Fﾗ\uFF9Fﾘ\uFF9Fﾙ\uFF9Fﾚ\uFF9Fﾛ\uFF9Fﾜ\uFF9Fﾝ\uFF9F",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "を\u309Aぁ\u309Aぃ\u309Aぅ\u309Aぇ\u309Aぉ\u309Aゃ\u309Aゅ\u309Aょ\u309Aっ\u309Aー\u309Aあ\u309Aい\u309Aう\u309Aえ\u309Aお\u309Aか\u309Aき\u309Aく\u309Aけ\u309Aこ\u309Aさ\u309Aし\u309Aす\u309Aせ\u309Aそ\u309Aた\u309Aち\u309Aつ\u309Aて\u309Aと\u309Aな\u309Aに\u309Aぬ\u309Aね\u309Aの\u309Aぱぴぷぺぽま\u309Aみ\u309Aむ\u309Aめ\u309Aも\u309Aや\u309Aゆ\u309Aよ\u309Aら\u309Aり\u309Aる\u309Aれ\u309Aろ\u309Aわ\u309Aん\u309A",
		},
		{
			name:    "With (HalfKanaToFull | KatakanaToHiragana) Small Kana Extension",
			input:   "𛄲𛅐𛅑𛅒𛅕𛅤𛅥𛅦𛅧",
			options: kana.HalfwidthToWide | kana.KatakanaToHiragana,
			expect:  "𛄲𛅐𛅑𛅒𛄲𛅐𛅑𛅒𛅧",
		},
		{
			name:    "With FullwidthToNarrow Latin-1 Printable",
			input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
			options: kana.FullwidthToNarrow,
			expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		},
		{
			name:    "With FullwidthToNarrow General Punctuation Printable",
			input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
			options: kana.FullwidthToNarrow,
			expect:  "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		},
		{
			name:    "With FullwidthToNarrow Some of Mathematical Operators",
			input:   "−∥",
			options: kana.FullwidthToNarrow,
			expect:  "−∥",
		},
		{
			name:    "With FullwidthToNarrow CJK Symbols and Punctuation",
			input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
			options: kana.FullwidthToNarrow,
			expect:  " 、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		},
		{
			name:    "With FullwidthToNarrow Fullwidth Forms",
			input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
			options: kana.FullwidthToNarrow,
			expect:  "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~⦅⦆¢£¬¯¦¥₩",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := kana.Convert(tc.input, tc.options)
			if diff := cmp.Diff(actual, tc.expect); diff != "" {
				t.Errorf("diff (-actual +expect): %s", diff)
			}
		})
	}
}
