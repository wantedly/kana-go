package nkf_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wantedly/kana-go/nkf"
)

type convertTestcase struct {
	name    string
	input   string
	options string
	expect  string
}

var testcases = []convertTestcase{
	{
		name:    "Base case ASCII Printable",
		input:   " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		options: "-w -W -m0 -x",
		expect:  " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
	},
	{
		name:    "Base case Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -x",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "Base case Latin-1 Formatting",
		input:   "\u00A0\u00AD",
		options: "-w -W -m0 -x",
		expect:  "\u00A0\u00AD",
	},
	{
		name:    "Base case General Punctuation Printable (nondegenerate only)",
		input:   "‐‒–—‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -x",
		expect:  "‐‒–—‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "Base case Minus Sign",
		input:   "−",
		options: "-w -W -m0 -x",
		expect:  "−",
	},
	{
		name:    "Base case CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -x",
		expect:  "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "Base case Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -x",
		expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
	},
	{
		name:    "Base case Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -x",
		expect:  "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
	},
	{
		name:    "Base case Fullwidth forms (nondegenerate only)",
		input:   "！＂＃＄％＆＇（）＊＋，．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￦",
		options: "-w -W -m0 -x",
		expect:  "！＂＃＄％＆＇（）＊＋，．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￦",
	},
	{
		name:    "Base case Halfwidth forms",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
		options: "-w -W -m0 -x",
		expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
	},
	{
		name:    "Base case degenerate",
		input:   "―∥－￠￡￢￣￤￥",
		options: "-w -W -m0 -x",
		expect:  "—‖−¢£¬‾¦¥",
	},
	{
		name:    "With -h Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -x -h",
		expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
	},
	{
		name:    "With -h Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -x -h",
		expect:  "゠ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔヵヶヷヸヹヺ・ーゝゞヿ",
	},
	{
		name:    "With -h Halfwidth Katakana",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
		options: "-w -W -m0 -x -h",
		expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
	},
	{
		name:    "With -h2 Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -x -h2",
		expect:  "ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴゕゖ\u3099\u309A゛゜ヽヾゟ",
	},
	{
		name:    "With -h2 Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -x -h2",
		expect:  "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
	},
	{
		name:    "With -h2 Halfwidth Katakana",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
		options: "-w -W -m0 -x -h2",
		expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
	},
	{
		name:    "With -h3 Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -x -h3",
		expect:  "ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴゕゖ\u3099\u309A゛゜ヽヾゟ",
	},
	{
		name:    "With -h3 Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -x -h3",
		expect:  "゠ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔヵヶヷヸヹヺ・ーゝゞヿ",
	},
	{
		name:    "With -h3 Halfwidth Katakana",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
		options: "-w -W -m0 -x -h3",
		expect:  "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
	},
	{
		name:    "With -X Basic Halfwidth Katakana",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
		options: "-w -W -m0",
		expect:  "。「」、・ヲァィゥェォャュョッーアイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワン゛゜",
	},
	{
		name:    "With -X Voiced Composites",
		input:   "ｦ\uFF9Eｧ\uFF9Eｨ\uFF9Eｩ\uFF9Eｪ\uFF9Eｫ\uFF9Eｬ\uFF9Eｭ\uFF9Eｮ\uFF9Eｯ\uFF9Eｰ\uFF9Eｱ\uFF9Eｲ\uFF9Eｳ\uFF9Eｴ\uFF9Eｵ\uFF9Eｶ\uFF9Eｷ\uFF9Eｸ\uFF9Eｹ\uFF9Eｺ\uFF9Eｻ\uFF9Eｼ\uFF9Eｽ\uFF9Eｾ\uFF9Eｿ\uFF9Eﾀ\uFF9Eﾁ\uFF9Eﾂ\uFF9Eﾃ\uFF9Eﾄ\uFF9Eﾅ\uFF9Eﾆ\uFF9Eﾇ\uFF9Eﾈ\uFF9Eﾉ\uFF9Eﾊ\uFF9Eﾋ\uFF9Eﾌ\uFF9Eﾍ\uFF9Eﾎ\uFF9Eﾏ\uFF9Eﾐ\uFF9Eﾑ\uFF9Eﾒ\uFF9Eﾓ\uFF9Eﾔ\uFF9Eﾕ\uFF9Eﾖ\uFF9Eﾗ\uFF9Eﾘ\uFF9Eﾙ\uFF9Eﾚ\uFF9Eﾛ\uFF9Eﾜ\uFF9Eﾝ\uFF9E",
		options: "-w -W -m0",
		expect:  "ヲ゛ァ゛ィ゛ゥ゛ェ゛ォ゛ャ゛ュ゛ョ゛ッ゛ー゛ア゛イ゛ヴエ゛オ゛ガギグゲゴザジズゼゾダヂヅデドナ゛ニ゛ヌ゛ネ゛ノ゛バビブベボマ゛ミ゛ム゛メ゛モ゛ヤ゛ユ゛ヨ゛ラ゛リ゛ル゛レ゛ロ゛ワ゛ン゛",
	},
	{
		name:    "With -X Semi-Voiced Composites",
		input:   "ｦ\uFF9Fｧ\uFF9Fｨ\uFF9Fｩ\uFF9Fｪ\uFF9Fｫ\uFF9Fｬ\uFF9Fｭ\uFF9Fｮ\uFF9Fｯ\uFF9Fｰ\uFF9Fｱ\uFF9Fｲ\uFF9Fｳ\uFF9Fｴ\uFF9Fｵ\uFF9Fｶ\uFF9Fｷ\uFF9Fｸ\uFF9Fｹ\uFF9Fｺ\uFF9Fｻ\uFF9Fｼ\uFF9Fｽ\uFF9Fｾ\uFF9Fｿ\uFF9Fﾀ\uFF9Fﾁ\uFF9Fﾂ\uFF9Fﾃ\uFF9Fﾄ\uFF9Fﾅ\uFF9Fﾆ\uFF9Fﾇ\uFF9Fﾈ\uFF9Fﾉ\uFF9Fﾊ\uFF9Fﾋ\uFF9Fﾌ\uFF9Fﾍ\uFF9Fﾎ\uFF9Fﾏ\uFF9Fﾐ\uFF9Fﾑ\uFF9Fﾒ\uFF9Fﾓ\uFF9Fﾔ\uFF9Fﾕ\uFF9Fﾖ\uFF9Fﾗ\uFF9Fﾘ\uFF9Fﾙ\uFF9Fﾚ\uFF9Fﾛ\uFF9Fﾜ\uFF9Fﾝ\uFF9F",
		options: "-w -W -m0",
		expect:  "ヲ゜ァ゜ィ゜ゥ゜ェ゜ォ゜ャ゜ュ゜ョ゜ッ゜ー゜ア゜イ゜ウ゜エ゜オ゜カ゜キ゜ク゜ケ゜コ゜サ゜シ゜ス゜セ゜ソ゜タ゜チ゜ツ゜テ゜ト゜ナ゜ニ゜ヌ゜ネ゜ノ゜パピプペポマ゜ミ゜ム゜メ゜モ゜ヤ゜ユ゜ヨ゜ラ゜リ゜ル゜レ゜ロ゜ワ゜ン゜",
	},
	{
		name:    "With -X -h Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -h",
		expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
	},
	{
		name:    "With -X -h Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -h",
		expect:  "゠ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔヵヶヷヸヹヺ・ーゝゞヿ",
	},
	{
		name:    "With -X -h Halfwidth Katakana",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F",
		options: "-w -W -m0 -h",
		expect:  "。「」、・をぁぃぅぇぉゃゅょっーあいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわん゛゜",
	},
	{
		name:    "With -X -h Voiced Composites",
		input:   "ｦ\uFF9Eｧ\uFF9Eｨ\uFF9Eｩ\uFF9Eｪ\uFF9Eｫ\uFF9Eｬ\uFF9Eｭ\uFF9Eｮ\uFF9Eｯ\uFF9Eｰ\uFF9Eｱ\uFF9Eｲ\uFF9Eｳ\uFF9Eｴ\uFF9Eｵ\uFF9Eｶ\uFF9Eｷ\uFF9Eｸ\uFF9Eｹ\uFF9Eｺ\uFF9Eｻ\uFF9Eｼ\uFF9Eｽ\uFF9Eｾ\uFF9Eｿ\uFF9Eﾀ\uFF9Eﾁ\uFF9Eﾂ\uFF9Eﾃ\uFF9Eﾄ\uFF9Eﾅ\uFF9Eﾆ\uFF9Eﾇ\uFF9Eﾈ\uFF9Eﾉ\uFF9Eﾊ\uFF9Eﾋ\uFF9Eﾌ\uFF9Eﾍ\uFF9Eﾎ\uFF9Eﾏ\uFF9Eﾐ\uFF9Eﾑ\uFF9Eﾒ\uFF9Eﾓ\uFF9Eﾔ\uFF9Eﾕ\uFF9Eﾖ\uFF9Eﾗ\uFF9Eﾘ\uFF9Eﾙ\uFF9Eﾚ\uFF9Eﾛ\uFF9Eﾜ\uFF9Eﾝ\uFF9E",
		options: "-w -W -m0 -h",
		expect:  "を゛ぁ゛ぃ゛ぅ゛ぇ゛ぉ゛ゃ゛ゅ゛ょ゛っ゛ー゛あ゛い゛ゔえ゛お゛がぎぐげござじずぜぞだぢづでどな゛に゛ぬ゛ね゛の゛ばびぶべぼま゛み゛む゛め゛も゛や゛ゆ゛よ゛ら゛り゛る゛れ゛ろ゛わ゛ん゛",
	},
	{
		name:    "With -X -h Semi-Voiced Composites",
		input:   "ｦ\uFF9Fｧ\uFF9Fｨ\uFF9Fｩ\uFF9Fｪ\uFF9Fｫ\uFF9Fｬ\uFF9Fｭ\uFF9Fｮ\uFF9Fｯ\uFF9Fｰ\uFF9Fｱ\uFF9Fｲ\uFF9Fｳ\uFF9Fｴ\uFF9Fｵ\uFF9Fｶ\uFF9Fｷ\uFF9Fｸ\uFF9Fｹ\uFF9Fｺ\uFF9Fｻ\uFF9Fｼ\uFF9Fｽ\uFF9Fｾ\uFF9Fｿ\uFF9Fﾀ\uFF9Fﾁ\uFF9Fﾂ\uFF9Fﾃ\uFF9Fﾄ\uFF9Fﾅ\uFF9Fﾆ\uFF9Fﾇ\uFF9Fﾈ\uFF9Fﾉ\uFF9Fﾊ\uFF9Fﾋ\uFF9Fﾌ\uFF9Fﾍ\uFF9Fﾎ\uFF9Fﾏ\uFF9Fﾐ\uFF9Fﾑ\uFF9Fﾒ\uFF9Fﾓ\uFF9Fﾔ\uFF9Fﾕ\uFF9Fﾖ\uFF9Fﾗ\uFF9Fﾘ\uFF9Fﾙ\uFF9Fﾚ\uFF9Fﾛ\uFF9Fﾜ\uFF9Fﾝ\uFF9F",
		options: "-w -W -m0 -h",
		expect:  "を゜ぁ゜ぃ゜ぅ゜ぇ゜ぉ゜ゃ゜ゅ゜ょ゜っ゜ー゜あ゜い゜う゜え゜お゜か゜き゜く゜け゜こ゜さ゜し゜す゜せ゜そ゜た゜ち゜つ゜て゜と゜な゜に゜ぬ゜ね゜の゜ぱぴぷぺぽま゜み゜む゜め゜も゜や゜ゆ゜よ゜ら゜り゜る゜れ゜ろ゜わ゜ん゜",
	},
	{
		name:    "With -Z Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -x -Z",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³'µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "With -Z General Punctuation Printable",
		input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -x -Z",
		expect:  "‐‒–--‖‗`'‚‛\"\"„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "With -Z Minus Sign",
		input:   "−",
		options: "-w -W -m0 -x -Z",
		expect:  "-",
	},
	{
		name:    "With -Z CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -x -Z",
		expect:  "　、。〃〄々〆〇<>《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "With -Z Fullwidth Forms",
		input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		options: "-w -W -m0 -x -Z",
		expect:  "!＂#$%&＇()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}～｟｠¢£¬‾¦¥￦",
	},
	{
		name:    "With -Z1 Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -x -Z1",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³'µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "With -Z1 General Punctuation Printable",
		input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -x -Z1",
		expect:  "‐‒–--‖‗`'‚‛\"\"„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "With -Z1 Minus Sign",
		input:   "−",
		options: "-w -W -m0 -x -Z1",
		expect:  "-",
	},
	{
		name:    "With -Z1 CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -x -Z1",
		expect:  " 、。〃〄々〆〇<>《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "With -Z1 Fullwidth Forms",
		input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		options: "-w -W -m0 -x -Z1",
		expect:  "!＂#$%&＇()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}～｟｠¢£¬‾¦¥￦",
	},
	{
		name:    "With -Z2 Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -x -Z2",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³'µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "With -Z2 General Punctuation Printable",
		input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -x -Z2",
		expect:  "‐‒–--‖‗`'‚‛\"\"„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "With -Z2 Minus Sign",
		input:   "−",
		options: "-w -W -m0 -x -Z2",
		expect:  "-",
	},
	{
		name:    "With -Z2 CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -x -Z2",
		expect:  "  、。〃〄々〆〇<>《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "With -Z2 Fullwidth Forms",
		input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		options: "-w -W -m0 -x -Z2",
		expect:  "!＂#$%&＇()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}～｟｠¢£¬‾¦¥￦",
	},
	{
		name:    "With -Z4 Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -x -Z4",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³'µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "With -Z4 General Punctuation Printable",
		input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -x -Z4",
		expect:  "‐‒–--‖‗`'‚‛\"\"„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "With -Z4 Minus Sign",
		input:   "−",
		options: "-w -W -m0 -x -Z4",
		expect:  "-",
	},
	{
		name:    "With -Z4 CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -x -Z4",
		expect:  "　､｡〃〄々〆〇<>《》｢｣『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "With -Z4 Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -x -Z4",
		expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\uFF9E\uFF9F\uFF9E\uFF9Fゝゞゟ",
	},
	{
		name:    "With -Z4 Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -x -Z4",
		expect:  "゠ｧｱｨｲｩｳｪｴｫｵｶｶﾞｷｷﾞｸｸﾞｹｹﾞｺｺﾞｻｻﾞｼｼﾞｽｽﾞｾｾﾞｿｿﾞﾀﾀﾞﾁﾁﾞｯﾂﾂﾞﾃﾃﾞﾄﾄﾞﾅﾆﾇﾈﾉﾊﾊﾞﾊﾟﾋﾋﾞﾋﾟﾌﾌﾞﾌﾟﾍﾍﾞﾍﾟﾎﾎﾞﾎﾟﾏﾐﾑﾒﾓｬﾔｭﾕｮﾖﾗﾘﾙﾚﾛヮﾜヰヱｦﾝｳﾞヵヶヷヸヹヺ･ｰヽヾヿ",
	},
	{
		name:    "With -Z4 Fullwidth Forms",
		input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		options: "-w -W -m0 -x -Z4",
		expect:  "!＂#$%&＇()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}～｟｠¢£¬‾¦¥￦",
	},
	{
		name:    "With -X -Z4 Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -Z4",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³'µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "With -X -Z4 General Punctuation Printable",
		input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -Z4",
		expect:  "‐‒–--‖‗`'‚‛\"\"„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "With -X -Z4 Minus Sign",
		input:   "−",
		options: "-w -W -m0 -Z4",
		expect:  "-",
	},
	{
		name:    "With -X -Z4 CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -Z4",
		expect:  "　､｡〃〄々〆〇<>《》｢｣『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "With -X -Z4 Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -Z4",
		expect:  "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\uFF9E\uFF9F\uFF9E\uFF9Fゝゞゟ",
	},
	{
		name:    "With -X -Z4 Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -Z4",
		expect:  "゠ｧｱｨｲｩｳｪｴｫｵｶｶﾞｷｷﾞｸｸﾞｹｹﾞｺｺﾞｻｻﾞｼｼﾞｽｽﾞｾｾﾞｿｿﾞﾀﾀﾞﾁﾁﾞｯﾂﾂﾞﾃﾃﾞﾄﾄﾞﾅﾆﾇﾈﾉﾊﾊﾞﾊﾟﾋﾋﾞﾋﾟﾌﾌﾞﾌﾟﾍﾍﾞﾍﾟﾎﾎﾞﾎﾟﾏﾐﾑﾒﾓｬﾔｭﾕｮﾖﾗﾘﾙﾚﾛヮﾜヰヱｦﾝｳﾞヵヶヷヸヹヺ･ｰヽヾヿ",
	},
	{
		name:    "With -X -Z4 Fullwidth Forms",
		input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		options: "-w -W -m0 -Z4",
		expect:  "!＂#$%&＇()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}～｟｠¢£¬‾¦¥￦",
	},
	{
		name:    "With -X -Z4 Halfwidth forms",
		input:   "｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ\uFF9E\uFF9F\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
		options: "-w -W -m0 -Z4",
		expect:  "。「」、・ヲァィゥェォャュョッーアイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワン゛゜\uFFA0ﾡﾢﾣﾤﾥﾦﾧﾨﾩﾪﾫﾬﾭﾮﾯﾰﾱﾲﾳﾴﾵﾶﾷﾸﾹﾺﾻﾼﾽﾾￂￃￄￅￆￇￊￋￌￍￎￏￒￓￔￕￖￗￚￛￜ￨￩￪￫￬￭￮",
	},
	{
		name:    "With -h2 -Z4 Latin-1 Printable",
		input:   "¡¢£¤¥¦§¨©ª«¬®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "¡¢£¤¥¦§¨©ª«¬®¯°±²³'µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ",
	},
	{
		name:    "With -h2 -Z4 General Punctuation Printable",
		input:   "‐‒–—―‖‗‘’‚‛“”„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "‐‒–--‖‗`'‚‛\"\"„‟†‡•‣․‥…‧‰‱′″‴‵‶‷‸‹›※‼‽‾‿⁀⁁⁂⁃⁄⁅⁆⁇⁈⁉⁊⁋⁌⁍⁎⁏⁐⁑⁒⁓⁔⁕⁖⁗⁘⁙⁚⁛⁜⁝⁞",
	},
	{
		name:    "With -h2 -Z4 Minus Sign",
		input:   "−",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "-",
	},
	{
		name:    "With -h2 -Z4 CJK Symbols and Punctuation",
		input:   "　、。〃〄々〆〇〈〉《》「」『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "　､｡〃〄々〆〇<>《》｢｣『』【】〒〓〔〕〖〗〘〙〚〛〜〝〞〟〠〡〢〣〤〥〦〧〨〩〪〭〫〬\u302E\u302F〰〱〲〳〴〵〶〷〸〹〺〻〼〽\u303E\u303F",
	},
	{
		name:    "With -h2 -Z4 Hiragana",
		input:   "ぁあぃいぅうぇえぉおかがきぎくぐけげこごさざしじすずせぜそぞただちぢっつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもゃやゅゆょよらりるれろゎわゐゑをんゔゕゖ\u3099\u309A゛゜ゝゞゟ",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴゕゖ\uFF9E\uFF9F\uFF9E\uFF9Fヽヾゟ",
	},
	{
		name:    "With -h2 -Z4 Katakana",
		input:   "゠ァアィイゥウェエォオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂッツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモャヤュユョヨラリルレロヮワヰヱヲンヴヵヶヷヸヹヺ・ーヽヾヿ",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "゠ｧｱｨｲｩｳｪｴｫｵｶｶﾞｷｷﾞｸｸﾞｹｹﾞｺｺﾞｻｻﾞｼｼﾞｽｽﾞｾｾﾞｿｿﾞﾀﾀﾞﾁﾁﾞｯﾂﾂﾞﾃﾃﾞﾄﾄﾞﾅﾆﾇﾈﾉﾊﾊﾞﾊﾟﾋﾋﾞﾋﾟﾌﾌﾞﾌﾟﾍﾍﾞﾍﾟﾎﾎﾞﾎﾟﾏﾐﾑﾒﾓｬﾔｭﾕｮﾖﾗﾘﾙﾚﾛヮﾜヰヱｦﾝｳﾞヵヶヷヸヹヺ･ｰヽヾヿ",
	},
	{
		name:    "With -h2 -Z4 Fullwidth Forms",
		input:   "！＂＃＄％＆＇（）＊＋，－．／０１２３４５６７８９：；＜＝＞？＠ＡＢＣＤＥＦＧＨＩＪＫＬＭＮＯＰＱＲＳＴＵＶＷＸＹＺ［＼］＾＿｀ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｛｜｝～｟｠￠￡￢￣￤￥￦",
		options: "-w -W -m0 -x -h2 -Z4",
		expect:  "!＂#$%&＇()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}～｟｠¢£¬‾¦¥￦",
	},
}

func TestConvert(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := nkf.Convert(tc.input, tc.options)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(actual, tc.expect); diff != "" {
				t.Errorf("diff (-actual +expect): %s", diff)
			}
		})
	}
}

func TestRealNKF(t *testing.T) {
	noNKF := os.Getenv("NO_REAL_NKF")
	if noNKF == "1" || noNKF == "true" || noNKF == "yes" {
		t.Skip("NO_REAL_NKF is set; skipping this test")
	}
	nkfPath := os.Getenv("REAL_NKF")
	if nkfPath == "" {
		nkfPath = "nkf"
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			command := exec.Command(nkfPath, tc.options)
			command.Stdin = strings.NewReader(tc.input)
			output, err := command.Output()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(string(output), tc.expect); diff != "" {
				t.Errorf("diff (-actual +expect): %s", diff)
			}
		})
	}
}
