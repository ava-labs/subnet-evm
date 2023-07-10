"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const mnemonic_1 = __importDefault(require("src/utils/mnemonic"));
const randomBytes = require("randombytes");
const mnemonic = mnemonic_1.default.getInstance();
const entropy = "9d7c99e77261acb88a5ed717f625d5d3ed5569e0f60429cc6eb9c4e91f48fb7c";
const langs = [
    "japanese",
    "spanish",
    "italian",
    "french",
    "korean",
    "czech",
    "portuguese",
    "chinese_traditional",
    "chinese_simplified"
];
const mnemnonics = [
    "ていし　みっか　せんせい　みっか　えいぶん　さつたば　かいよう　へんたい　うやまう　にちじょう　せっかく　とける　ぶどう　にんぷ　たいうん　はいそう　かえる　したて　なめらか　だじゃれ　だんれつ　てんぼうだい　もくようび　そむく",
    "nueve tenis lágrima tenis baile folleto canica sonoro autor perla jardín oxígeno sensor piscina lonja rabo cañón germen pegar marrón molino opuesto trébol llaga",
    "pergamena tensione maratona tensione becco geco cena srotolato badilata regola lumaca prelievo somma rifasare motivato sarcasmo ceramica ibernato randagio ninfa orafo polmonite tuffato modulo",
    "mobile surface héron surface batterie éthanol capsule serein bafouer pangolin gravir nuisible salive peintre intense préfixe carabine fatal orque label lucide neurone toucher informer",
    "운반 특별 시인 특별 귀신 빗물 농민 취업 구입 작년 스님 이윽고 체험 장애인 아흔 제작 농장 상추 입사 언덕 염려 의외로 학급 씨름",
    "pohnutka vize nikam vize datum ledvina export uklidnit cirkus revolver naslepo procento traverza rozjezd odliv slavnost fajfka lyra rande omluva panovat poukaz vyrvat ochladit",
    "mesada surdina guincho surdina aumentar escrita brilho sediado assador ostentar goela nevoeiro rouco panqueca inovador postal brochura faceta ontem judoca linhagem munido torque indeciso",
    "烯 逮 岩 逮 資 衛 走 賦 料 默 膠 辛 杯 挑 戶 陶 議 劉 拍 謀 浮 休 煩 慮",
    "烯 逮 岩 逮 资 卫 走 赋 料 默 胶 辛 杯 挑 户 陶 议 刘 拍 谋 浮 休 烦 虑"
];
const seeds = [
    "2ed50c99aa2ee350f0a46c8427d69f9f5c7c5864be7a64ae96695374a0a5a02a3c5075312234f05f8f4c840aa486c99131f84b81c56daff5c31a89cdc7b50424",
    "04c6cfd9337642f47e21e28632f9744fd1b56c57454ebae5c627c2a4b16e47c0948b9c582041bbb1590e128a25ae79d7055ed8aecdbc030920a67205da24846d",
    "c4274544eb6c375d2caa70af8c6a67e3b751c518cbb35244891c7b74a12a5e06d5ce5b925f134930e17f5e86b21146d096715c7645a250dbba1d2ba35bc07317",
    "00e5b5e4785763d6f92fe1ad7c5a7e7dd0fd375bd441473198d2486990ca5a924b5cb6b6831dc94d446c9b3180eefe2d799887bcfc1ee6d8f4f0650e594c9d1b",
    "d8dcc049712867ba7d1bc0e2874d8ec38ee26088d1e2affa10ffac30cf1d0b915bbb6c79bfafbb1db0e8cd66880bf4ba52c53b949f6a3adbba1821dd3045c7cb",
    "a81d8a917861cb8a1ffd2e94452e08fd6dc4d2577bad3ac089f56279521b1c95caebe57ac6c3d126d8abd4d6a3f2c3d8c207bd36fbf75a5e797c5a8f1992c651",
    "b9fc39d7f138a95b8f31436e02a8711b3164cd45a211673f7137429b45faf77a1604682b51803a983638c46a8b2c13237c87ab4b685a1fa5c65700dc7136ccfc",
    "1a5f3eab01276bf7d9b06c42be90fb4b8106b278b4bbaf922f3da6821a63b4d69fa9285fec0936e0f04a1b2a25a65064cd51c111c75063dbe59e4de336b35f85",
    "53bcb9fe403a4a45bee2a2a04dabfa3f2018db349d0e5200175bd345aaa3a3bd45a88a6fb7244914ad156961742f7b4ea7f8ffd83fcae5b1166b73b2ad943f76"
];
const vectors = [
    {
        entropy: "00000000000000000000000000000000",
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
        seed: "c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04"
    },
    {
        entropy: "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
        mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow",
        seed: "2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607"
    },
    {
        entropy: "80808080808080808080808080808080",
        mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
        seed: "d71de856f81a8acc65e6fc851a38d4d7ec216fd0796d0a6827a3ad6ed5511a30fa280f12eb2e47ed2ac03b5c462a0358d18d69fe4f985ec81778c1b370b652a8"
    },
    {
        entropy: "ffffffffffffffffffffffffffffffff",
        mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
        seed: "ac27495480225222079d7be181583751e86f571027b0497b5b5d11218e0a8a13332572917f0f8e5a589620c6f15b11c61dee327651a14c34e18231052e48c069"
    },
    {
        entropy: "000000000000000000000000000000000000000000000000",
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon agent",
        seed: "035895f2f481b1b0f01fcf8c289c794660b289981a78f8106447707fdd9666ca06da5a9a565181599b79f53b844d8a71dd9f439c52a3d7b3e8a79c906ac845fa"
    },
    {
        entropy: "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
        mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
        seed: "f2b94508732bcbacbcc020faefecfc89feafa6649a5491b8c952cede496c214a0c7b3c392d168748f2d4a612bada0753b52a1c7ac53c1e93abd5c6320b9e95dd"
    },
    {
        entropy: "808080808080808080808080808080808080808080808080",
        mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
        seed: "107d7c02a5aa6f38c58083ff74f04c607c2d2c0ecc55501dadd72d025b751bc27fe913ffb796f841c49b1d33b610cf0e91d3aa239027f5e99fe4ce9e5088cd65"
    },
    {
        entropy: "ffffffffffffffffffffffffffffffffffffffffffffffff",
        mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
        seed: "0cd6e5d827bb62eb8fc1e262254223817fd068a74b5b449cc2f667c3f1f985a76379b43348d952e2265b4cd129090758b3e3c2c49103b5051aac2eaeb890a528"
    },
    {
        entropy: "0000000000000000000000000000000000000000000000000000000000000000",
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
        seed: "bda85446c68413707090a52022edd26a1c9462295029f2e60cd7c4f2bbd3097170af7a4d73245cafa9c3cca8d561a7c3de6f5d4a10be8ed2a5e608d68f92fcc8"
    },
    {
        entropy: "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
        mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title",
        seed: "bc09fca1804f7e69da93c2f2028eb238c227f2e9dda30cd63699232578480a4021b146ad717fbb7e451ce9eb835f43620bf5c514db0f8add49f5d121449d3e87"
    },
    {
        entropy: "8080808080808080808080808080808080808080808080808080808080808080",
        mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
        seed: "c0c519bd0e91a2ed54357d9d1ebef6f5af218a153624cf4f2da911a0ed8f7a09e2ef61af0aca007096df430022f7a2b6fb91661a9589097069720d015e4e982f"
    },
    {
        entropy: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
        mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
        seed: "dd48c104698c30cfe2b6142103248622fb7bb0ff692eebb00089b32d22484e1613912f0a5b694407be899ffd31ed3992c456cdf60f5d4564b8ba3f05a69890ad"
    },
    {
        entropy: "77c2b00716cec7213839159e404db50d",
        mnemonic: "jelly better achieve collect unaware mountain thought cargo oxygen act hood bridge",
        seed: "b5b6d0127db1a9d2226af0c3346031d77af31e918dba64287a1b44b8ebf63cdd52676f672a290aae502472cf2d602c051f3e6f18055e84e4c43897fc4e51a6ff"
    },
    {
        entropy: "b63a9c59a6e641f288ebc103017f1da9f8290b3da6bdef7b",
        mnemonic: "renew stay biology evidence goat welcome casual join adapt armor shuffle fault little machine walk stumble urge swap",
        seed: "9248d83e06f4cd98debf5b6f010542760df925ce46cf38a1bdb4e4de7d21f5c39366941c69e1bdbf2966e0f6e6dbece898a0e2f0a4c2b3e640953dfe8b7bbdc5"
    },
    {
        entropy: "3e141609b97933b66a060dcddc71fad1d91677db872031e85f4c015c5e7e8982",
        mnemonic: "dignity pass list indicate nasty swamp pool script soccer toe leaf photo multiply desk host tomato cradle drill spread actor shine dismiss champion exotic",
        seed: "ff7f3184df8696d8bef94b6c03114dbee0ef89ff938712301d27ed8336ca89ef9635da20af07d4175f2bf5f3de130f39c9d9e8dd0472489c19b1a020a940da67"
    },
    {
        entropy: "0460ef47585604c5660618db2e6a7e7f",
        mnemonic: "afford alter spike radar gate glance object seek swamp infant panel yellow",
        seed: "65f93a9f36b6c85cbe634ffc1f99f2b82cbb10b31edc7f087b4f6cb9e976e9faf76ff41f8f27c99afdf38f7a303ba1136ee48a4c1e7fcd3dba7aa876113a36e4"
    },
    {
        entropy: "72f60ebac5dd8add8d2a25a797102c3ce21bc029c200076f",
        mnemonic: "indicate race push merry suffer human cruise dwarf pole review arch keep canvas theme poem divorce alter left",
        seed: "3bbf9daa0dfad8229786ace5ddb4e00fa98a044ae4c4975ffd5e094dba9e0bb289349dbe2091761f30f382d4e35c4a670ee8ab50758d2c55881be69e327117ba"
    },
    {
        entropy: "2c85efc7f24ee4573d2b81a6ec66cee209b2dcbd09d8eddc51e0215b0b68e416",
        mnemonic: "clutch control vehicle tonight unusual clog visa ice plunge glimpse recipe series open hour vintage deposit universe tip job dress radar refuse motion taste",
        seed: "fe908f96f46668b2d5b37d82f558c77ed0d69dd0e7e043a5b0511c48c2f1064694a956f86360c93dd04052a8899497ce9e985ebe0c8c52b955e6ae86d4ff4449"
    },
    {
        entropy: "eaebabb2383351fd31d703840b32e9e2",
        mnemonic: "turtle front uncle idea crush write shrug there lottery flower risk shell",
        seed: "bdfb76a0759f301b0b899a1e3985227e53b3f51e67e3f2a65363caedf3e32fde42a66c404f18d7b05818c95ef3ca1e5146646856c461c073169467511680876c"
    },
    {
        entropy: "7ac45cfe7722ee6c7ba84fbc2d5bd61b45cb2fe5eb65aa78",
        mnemonic: "kiss carry display unusual confirm curtain upgrade antique rotate hello void custom frequent obey nut hole price segment",
        seed: "ed56ff6c833c07982eb7119a8f48fd363c4a9b1601cd2de736b01045c5eb8ab4f57b079403485d1c4924f0790dc10a971763337cb9f9c62226f64fff26397c79"
    },
    {
        entropy: "4fa1a8bc3e6d80ee1316050e862c1812031493212b7ec3f3bb1b08f168cabeef",
        mnemonic: "exile ask congress lamp submit jacket era scheme attend cousin alcohol catch course end lucky hurt sentence oven short ball bird grab wing top",
        seed: "095ee6f817b4c2cb30a5a797360a81a40ab0f9a4e25ecd672a3f58a0b5ba0687c096a6b14d2c0deb3bdefce4f61d01ae07417d502429352e27695163f7447a8c"
    },
    {
        entropy: "18ab19a9f54a9274f03e5209a2ac8a91",
        mnemonic: "board flee heavy tunnel powder denial science ski answer betray cargo cat",
        seed: "6eff1bb21562918509c73cb990260db07c0ce34ff0e3cc4a8cb3276129fbcb300bddfe005831350efd633909f476c45c88253276d9fd0df6ef48609e8bb7dca8"
    },
    {
        entropy: "18a2e1d81b8ecfb2a333adcb0c17a5b9eb76cc5d05db91a4",
        mnemonic: "board blade invite damage undo sun mimic interest slam gaze truly inherit resist great inject rocket museum chief",
        seed: "f84521c777a13b61564234bf8f8b62b3afce27fc4062b51bb5e62bdfecb23864ee6ecf07c1d5a97c0834307c5c852d8ceb88e7c97923c0a3b496bedd4e5f88a9"
    },
    {
        entropy: "15da872c95a13dd738fbf50e427583ad61f18fd99f628c417a61cf8343c90419",
        mnemonic: "beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut",
        seed: "b15509eaa2d09d3efd3e006ef42151b30367dc6e3aa5e44caba3fe4d3e352e65101fbdb86a96776b91946ff06f8eac594dc6ee1d3e82a42dfe1b40fef6bcc3fd"
    }
];
const badMnemonics = [
    {
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon"
    },
    {
        mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow yellow"
    },
    {
        mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice caged above"
    },
    { mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo, wrong" },
    {
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon"
    },
    {
        mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will will will"
    },
    {
        mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always."
    },
    {
        mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo why"
    },
    {
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art art"
    },
    {
        mnemonic: "legal winner thank year wave sausage worth useful legal winner thanks year wave worth useful legal winner thank year wave sausage worth title"
    },
    {
        mnemonic: "letter advice cage absurd amount doctor acoustic avoid letters advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless"
    },
    {
        mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo voted"
    },
    {
        mnemonic: "jello better achieve collect unaware mountain thought cargo oxygen act hood bridge"
    },
    {
        mnemonic: "renew, stay, biology, evidence, goat, welcome, casual, join, adapt, armor, shuffle, fault, little, machine, walk, stumble, urge, swap"
    },
    { mnemonic: "dignity pass list indicate nasty" },
    // From issue 32
    {
        mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon letter"
    }
];
const malformedMnemonics = [
    "a a a a a a a a a a a a a a a a a a a a a a a a a",
    "a",
    "a a a a a a a a a a a a a a" // Not multiple of 3
];
describe("Mnemonic", () => {
    test("vectors", () => {
        vectors.forEach((vector, index) => __awaiter(void 0, void 0, void 0, function* () {
            const wordlist = mnemonic.getWordlists("english");
            const entropyToMnemonic = mnemonic.entropyToMnemonic(vector.entropy, wordlist);
            expect(vector.mnemonic).toBe(entropyToMnemonic);
            const mnemonicToEntropy = mnemonic.mnemonicToEntropy(vector.mnemonic, wordlist);
            expect(mnemonicToEntropy).toBe(vector.entropy);
            const password = "TREZOR";
            const mnemonicToSeed = yield mnemonic.mnemonicToSeed(vector.mnemonic, password);
            expect(mnemonicToSeed.toString("hex")).toBe(vector.seed);
        }));
    });
    test("badMnemonics", () => {
        const wordlist = mnemonic.getWordlists("english");
        badMnemonics.forEach((badMnemonic, index) => {
            const validateMnemonic = mnemonic.validateMnemonic(badMnemonic.mnemonic, wordlist);
            expect(validateMnemonic).toBeFalsy();
        });
    });
    test("malformedMnemonics", () => {
        const wordlist = mnemonic.getWordlists("english");
        malformedMnemonics.forEach((malformedMnemonic, index) => {
            const validateMnemonic = mnemonic.validateMnemonic(malformedMnemonic, wordlist);
            expect(validateMnemonic).toBeFalsy();
        });
    });
    test("entropyToMnemonic", () => {
        langs.forEach((lang, index) => {
            const wordlist = mnemonic.getWordlists(lang);
            const entropyToMnemonic = mnemonic.entropyToMnemonic(entropy, wordlist);
            expect(mnemnonics[index]).toBe(entropyToMnemonic);
        });
    });
    test("generateMnemonic", () => {
        const strength = 256;
        langs.forEach((lang) => {
            const wordlist = mnemonic.getWordlists(lang);
            const m = mnemonic.generateMnemonic(strength, randomBytes, wordlist);
            expect(typeof m === "string").toBeTruthy();
        });
    });
    test("test mnemonic lengths", () => {
        const wordlist = mnemonic.getWordlists("english");
        let strength = 128;
        let mnemnnic = mnemonic.generateMnemonic(strength, randomBytes, wordlist);
        expect(mnemnnic.split(" ").length).toBe(12);
        strength = 160;
        mnemnnic = mnemonic.generateMnemonic(strength, randomBytes, wordlist);
        expect(mnemnnic.split(" ").length).toBe(15);
        strength = 192;
        mnemnnic = mnemonic.generateMnemonic(strength, randomBytes, wordlist);
        expect(mnemnnic.split(" ").length).toBe(18);
        strength = 224;
        mnemnnic = mnemonic.generateMnemonic(strength, randomBytes, wordlist);
        expect(mnemnnic.split(" ").length).toBe(21);
        strength = 256;
        mnemnnic = mnemonic.generateMnemonic(strength, randomBytes, wordlist);
        expect(mnemnnic.split(" ").length).toBe(24);
    });
    test("getWordlists", () => {
        langs.forEach((lang) => {
            const wordlist = mnemonic.getWordlists(lang);
            expect(typeof wordlist === "object").toBeTruthy();
        });
    });
    test("mnemonicToEntropy", () => {
        mnemnonics.forEach((mnemnnic, index) => {
            const wordlist = mnemonic.getWordlists(langs[index]);
            const mnemonicToEntropy = mnemonic.mnemonicToEntropy(mnemnnic, wordlist);
            expect(mnemonicToEntropy).toBe(entropy);
        });
    });
    test("mnemonicToSeed", () => __awaiter(void 0, void 0, void 0, function* () {
        mnemnonics.forEach((mnemnnic) => __awaiter(void 0, void 0, void 0, function* () {
            const password = "password";
            const mnemonicToSeed = yield mnemonic.mnemonicToSeed(mnemnnic, password);
            expect(typeof mnemonicToSeed === "object").toBeTruthy();
        }));
    }));
    test("mnemonicToSeedSync", () => {
        mnemnonics.forEach((mnemnnic, index) => {
            const password = "password";
            const mnemonicToSeedSync = mnemonic.mnemonicToSeedSync(mnemnnic, password);
            expect(mnemonicToSeedSync.toString("hex")).toBe(seeds[index]);
        });
    });
    test("validateMnemonic", () => {
        mnemnonics.forEach((mnemnnic, index) => {
            const wordlist = mnemonic.getWordlists(langs[index]);
            const validateMnemonic = mnemonic.validateMnemonic(mnemnnic, wordlist);
            expect(validateMnemonic).toBeTruthy();
        });
    });
    test("setDefaultWordlist", () => {
        langs.forEach((lang, index) => {
            mnemonic.setDefaultWordlist(lang);
            const getDefaultWordlist = mnemonic.getDefaultWordlist();
            const wordlist = mnemonic.getWordlists(lang);
            const m = mnemonic.generateMnemonic(256, randomBytes, wordlist);
            expect(getDefaultWordlist).toBe(lang);
            expect(typeof wordlist === "object").toBeTruthy();
            expect(typeof m === "string").toBeTruthy();
        });
    });
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibW5lbW9uaWMudGVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3Rlc3RzL3V0aWxzL21uZW1vbmljLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7QUFBQSxrRUFBeUM7QUFHekMsTUFBTSxXQUFXLEdBQVEsT0FBTyxDQUFDLGFBQWEsQ0FBQyxDQUFBO0FBQy9DLE1BQU0sUUFBUSxHQUFHLGtCQUFRLENBQUMsV0FBVyxFQUFFLENBQUE7QUFDdkMsTUFBTSxPQUFPLEdBQ1gsa0VBQWtFLENBQUE7QUFDcEUsTUFBTSxLQUFLLEdBQWE7SUFDdEIsVUFBVTtJQUNWLFNBQVM7SUFDVCxTQUFTO0lBQ1QsUUFBUTtJQUNSLFFBQVE7SUFDUixPQUFPO0lBQ1AsWUFBWTtJQUNaLHFCQUFxQjtJQUNyQixvQkFBb0I7Q0FDckIsQ0FBQTtBQUVELE1BQU0sVUFBVSxHQUFhO0lBQzNCLGdJQUFnSTtJQUNoSSx5S0FBeUs7SUFDekssaU1BQWlNO0lBQ2pNLDRMQUE0TDtJQUM1TCxnS0FBZ0s7SUFDaEssaUxBQWlMO0lBQ2pMLDRMQUE0TDtJQUM1TCxpREFBaUQ7SUFDakQsaURBQWlEO0NBQ2xELENBQUE7QUFFRCxNQUFNLEtBQUssR0FBYTtJQUN0QixrSUFBa0k7SUFDbEksa0lBQWtJO0lBQ2xJLGtJQUFrSTtJQUNsSSxrSUFBa0k7SUFDbEksa0lBQWtJO0lBQ2xJLGtJQUFrSTtJQUNsSSxrSUFBa0k7SUFDbEksa0lBQWtJO0lBQ2xJLGtJQUFrSTtDQUNuSSxDQUFBO0FBUUQsTUFBTSxPQUFPLEdBQWE7SUFDeEI7UUFDRSxPQUFPLEVBQUUsa0NBQWtDO1FBQzNDLFFBQVEsRUFDTiwrRkFBK0Y7UUFDakcsSUFBSSxFQUFFLGtJQUFrSTtLQUN6STtJQUNEO1FBQ0UsT0FBTyxFQUFFLGtDQUFrQztRQUMzQyxRQUFRLEVBQ04sNkVBQTZFO1FBQy9FLElBQUksRUFBRSxrSUFBa0k7S0FDekk7SUFDRDtRQUNFLE9BQU8s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LFlBQVksQ0FBQyxTQUFTLENBQWEsQ0FBQTtRQUM3RCxJQUFJLFFBQVEsR0FBVyxHQUFHLENBQUE7UUFDMUIsSUFBSSxRQUFRLEdBQVcsUUFBUSxDQUFDLGdCQUFnQixDQUM5QyxRQUFRLEVBQ1IsV0FBVyxFQUNYLFFBQVEsQ0FDVCxDQUFBO1FBQ0QsTUFBTSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsR0FBRyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQzNDLFFBQVEsR0FBRyxHQUFHLENBQUE7UUFDZCxRQUFRLEdBQUcsUUFBUSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsRUFBRSxXQUFXLEVBQUUsUUFBUSxDQUFDLENBQUE7UUFDckUsTUFBTSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsR0FBRyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQzNDLFFBQVEsR0FBRyxHQUFHLENBQUE7UUFDZCxRQUFRLEdBQUcsUUFBUSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsRUFBRSxXQUFXLEVBQUUsUUFBUSxDQUFDLENBQUE7UUFDckUsTUFBTSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsR0FBRyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQzNDLFFBQVEsR0FBRyxHQUFHLENBQUE7UUFDZCxRQUFRLEdBQUcsUUFBUSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsRUFBRSxXQUFXLEVBQUUsUUFBUSxDQUFDLENBQUE7UUFDckUsTUFBTSxDQUFDLFFBQVEsQ0FBQyxLQUFLLENBQUMsR0FBRyxDQUFDLENBQUMsTUFBTSxDQUFDLENBQUMsSUFBSSxDQUFDLEVBQUUsQ0FBQyxDQUFBO1FBQzNDLFFBQVEsR0FBRyxHQUFHLENBQUE7UUFDZCxRQUFRLEdBQUcsUUFBUSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsRUFBRSxXQUFXLEVBQUUsUUFBUSxDQUFDLENBQUE7UUFDckUsTUFB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QUMsa0JBQWtCLEVBQUUsQ0FBQTtZQUNoRSxNQUFNLFFBQVEsR0FBRyxRQUFRLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBYSxDQUFBO1lBQ3hELE1BQU0sQ0FBQyxHQUFXLFFBQVEsQ0FBQyxnQkFBZ0IsQ0FBQyxHQUFHLEVBQUUsV0FBVyxFQUFFLFFBQVEsQ0FBQyxDQUFBO1lBQ3ZFLE1BQU0sQ0FBQyxrQkFBa0IsQ0FBQyxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsQ0FBQTtZQUNyQyxNQUFNLENBQUMsT0FBTyxRQUFRLEtBQUssUUFBUSxDQUFDLENBQUMsVUFBVSxFQUFFLENBQUE7WUFDakQsTUFBTSxDQUFDLE9BQU8sQ0FBQyxLQUFLLFFBQVEsQ0FBQyxDQUFDLFVBQVUsRUFBRSxDQUFBO1FBQzVDLENBQUMsQ0FBQyxDQUFBO0lBQ0osQ0FBQyxDQUFDLENBQUE7QUFDSixDQUFDLENBQUMsQ0FBQSIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBNbmVtb25pYyBmcm9tIFwic3JjL3V0aWxzL21uZW1vbmljXCJcbmltcG9ydCB7IEJ1ZmZlciB9IGZyb20gXCJidWZmZXIvXCJcblxuY29uc3QgcmFuZG9tQnl0ZXM6IGFueSA9IHJlcXVpcmUoXCJyYW5kb21ieXRlc1wiKVxuY29uc3QgbW5lbW9uaWMgPSBNbmVtb25pYy5nZXRJbnN0YW5jZSgpXG5jb25zdCBlbnRyb3B5OiBzdHJpbmcgPVxuICBcIjlkN2M5OWU3NzI2MWFjYjg4YTVlZDcxN2Y2MjVkNWQzZWQ1NTY5ZTBmNjA0MjljYzZlYjljNGU5MWY0OGZiN2NcIlxuY29uc3QgbGFuZ3M6IHN0cmluZ1tdID0gW1xuICBcImphcGFuZXNlXCIsXG4gIFwic3BhbmlzaFwiLFxuICBcIml0YWxpYW5cIixcbiAgXCJmcmVuY2hcIixcbiAgXCJrb3JlYW5cIixcbiAgXCJjemVjaFwiLFxuICBcInBvcnR1Z3Vlc2VcIixcbiAgXCJjaGluZXNlX3RyYWRpdGlvbmFsXCIsXG4gIFwiY2hpbmVzZV9zaW1wbGlmaWVkXCJcbl1cblxuY29uc3QgbW5lbW5vbmljczogc3RyaW5nW10gPSBbXG4gIFwi44Gm44GE44GX44CA44G/44Gj44GL44CA44Gb44KT44Gb44GE44CA44G/44Gj44GL44CA44GI44GE44G144KZ44KT44CA44GV44Gk44Gf44Gv44KZ44CA44GL44GE44KI44GG44CA44G444KT44Gf44GE44CA44GG44KE44G+44GG44CA44Gr44Gh44GX44KZ44KH44GG44CA44Gb44Gj44GL44GP44CA44Go44GR44KL44CA44G144KZ44Go44KZ44GG44CA44Gr44KT44G144Ka44CA44Gf44GE44GG44KT44CA44Gv44GE44Gd44GG44CA44GL44GI44KL44CA44GX44Gf44Gm44CA44Gq44KB44KJ44GL44CA44Gf44KZ44GX44KZ44KD44KM44CA44Gf44KZ44KT44KM44Gk44CA44Gm44KT44G744KZ44GG44Gf44KZ44GE44CA44KC44GP44KI44GG44Gy44KZ44CA44Gd44KA44GPXCIsXG4gIFwibnVldmUgdGVuaXMgbGHMgWdyaW1hIHRlbmlzIGJhaWxlIGZvbGxldG8gY2FuaWNhIHNvbm9ybyBhdXRvciBwZXJsYSBqYXJkacyBbiBveGnMgWdlbm8gc2Vuc29yIHBpc2NpbmEgbG9uamEgcmFibyBjYW7Mg2/MgW4gZ2VybWVuIHBlZ2FyIG1hcnJvzIFuIG1vbGlubyBvcHVlc3RvIHRyZcyBYm9sIGxsYWdhXCIsXG4gIFwicGVyZ2FtZW5hIHRlbnNpb25lIG1hcmF0b25hIHRlbnNpb25lIGJlY2NvIGdlY28gY2VuYSBzcm90b2xhdG8gYmFkaWxhdGEgcmVnb2xhIGx1bWFjYSBwcmVsaWV2byBzb21tYSByaWZhc2FyZSBtb3RpdmF0byBzYXJjYXNtbyBjZXJhbWljYSBpYmVybmF0byByYW5kYWdpbyBuaW5mYSBvcmFmbyBwb2xtb25pdGUgdHVmZmF0byBtb2R1bG9cIixcbiAgXCJtb2JpbGUgc3VyZmFjZSBoZcyBcm9uIHN1cmZhY2UgYmF0dGVyaWUgZcyBdGhhbm9sIGNhcHN1bGUgc2VyZWluIGJhZm91ZXIgcGFuZ29saW4gZ3JhdmlyIG51aXNpYmxlIHNhbGl2ZSBwZWludHJlIGludGVuc2UgcHJlzIFmaXhlIGNhcmFiaW5lIGZhdGFsIG9ycXVlIGxhYmVsIGx1Y2lkZSBuZXVyb25lIHRvdWNoZXIgaW5mb3JtZXJcIixcbiAgXCLhhIvhha7hhqvhhIfhhaHhhqsg4YSQ4YWz4Yao4YSH4YWn4YavIOGEieGFteGEi+GFteGGqyDhhJDhhbPhhqjhhIfhhafhhq8g4YSA4YWx4YSJ4YW14YarIOGEh+GFteGGuuGEhuGFruGGryDhhILhhanhhrzhhIbhhbXhhqsg4YSO4YWx4YSL4YWl4Ya4IOGEgOGFruGEi+GFteGGuCDhhIzhhaHhhqjhhILhhafhhqsg4YSJ4YWz4YSC4YW14Ya3IOGEi+GFteGEi+GFs+GGqOGEgOGFqSDhhI7hhabhhJLhhaXhhrcg4YSM4YWh4Ya84YSL4YWi4YSL4YW14YarIOGEi+GFoeGEkuGFs+GGqyDhhIzhhabhhIzhhaHhhqgg4YSC4YWp4Ya84YSM4YWh4Ya8IOGEieGFoeGGvOGEjuGFriDhhIvhhbXhhrjhhInhhaEg4YSL4YWl4Yar4YSD4YWl4YaoIOGEi+GFp+GGt+GEheGFpyDhhIvhhbThhIvhhazhhIXhhakg4YSS4YWh4Yao4YSA4YWz4Ya4IOGEiuGFteGEheGFs+GGt1wiLFxuICBcInBvaG51dGthIHZpemUgbmlrYW0gdml6ZSBkYXR1bSBsZWR2aW5hIGV4cG9ydCB1a2xpZG5pdCBjaXJrdXMgcmV2b2x2ZXIgbmFzbGVwbyBwcm9jZW50byB0cmF2ZXJ6YSByb3pqZXpkIG9kbGl2IHNsYXZub3N0IGZhamZrYSBseXJhIHJhbmRlIG9tbHV2YSBwYW5vdmF0IHBvdWtheiB2eXJ2YXQgb2NobGFkaXRcIixcbiAgXCJtZXNhZGEgc3VyZGluYSBndWluY2hvIHN1cmRpbmEgYXVtZW50YXIgZXNjcml0YSBicmlsaG8gc2VkaWFkbyBhc3NhZG9yIG9zdGVudGFyIGdvZWxhIG5ldm9laXJvIHJvdWNvIHBhbnF1ZWNhIGlub3ZhZG9yIHBvc3RhbCBicm9jaHVyYSBmYWNldGEgb250ZW0ganVkb2NhIGxpbmhhZ2VtIG11bmlkbyB0b3JxdWUgaW5kZWNpc29cIixcbiAgXCLng68g6YCuIOWyqSDpgK4g6LOHIOihmyDotbAg6LOmIOaWmSDpu5gg6IagIOi+myDmna8g5oyRIOaItiDpmbYg6K2wIOWKiSDmi40g6KyAIOa1riDkvJEg54WpIOaFrlwiLFxuICBcIueDryDpgK4g5bKpIOmAriDotYQg5Y2rIOi1sCDotYsg5paZIOm7mCDog7Yg6L6bIOadryDmjJEg5oi3IOmZtiDorq4g5YiYIOaLjSDosIsg5rWuIOS8kSDng6Yg6JmRXCJcbl1cblxuY29uc3Qgc2VlZHM6IHN0cmluZ1tdID0gW1xuICBcIjJlZDUwYzk5YWEyZWUzNTBmMGE0NmM4NDI3ZDY5ZjlmNWM3YzU4NjRiZTdhNjRhZTk2Njk1Mzc0YTBhNWEwMmEzYzUwNzUzMTIyMzRmMDVmOGY0Yzg0MGFhNDg2Yzk5MTMxZjg0YjgxYzU2ZGFmZjVjMzFhODljZGM3YjUwNDI0XCIsXG4gIFwiMDRjNmNmZDkzMzc2NDJmNDdlMjFlMjg2MzJmOTc0NGZkMWI1NmM1NzQ1NGViYWU1YzYyN2MyYTRiMTZlNDdjMDk0OGI5YzU4MjA0MWJiYjE1OTBlMTI4YTI1YWU3OWQ3MDU1ZWQ4YWVjZGJjMDMwOTIwYTY3MjA1ZGEyNDg0NmRcIixcbiAgXCJjNDI3NDU0NGViNmMzNzVkMmNhYTcwYWY4YzZhNjdlM2I3NTFjNTE4Y2JiMzUyNDQ4OTFjN2I3NGExMmE1ZTA2ZDVjZTViOTI1ZjEzNDkzMGUxN2Y1ZTg2YjIxMTQ2ZDA5NjcxNWM3NjQ1YTI1MGRiYmExZDJiYTM1YmMwNzMxN1wiLFxuICBcIjAwZTViNWU0Nzg1NzYzZDZmOTJmZTFhZDdjNWE3ZTdkZDBmZDM3NWJkNDQxNDczMTk4ZDI0ODY5OTBjYTVhOTI0YjVjYjZiNjgzMWRjOTRkNDQ2YzliMzE4MGVlZmUyZDc5OTg4N2JjZmMxZWU2ZDhmNGYwNjUwZTU5NGM5ZDFiXCIsXG4gIFwiZDhkY2MwNDk3MTI4NjdiYTdkMWJjMGUyODc0ZDhlYzM4ZWUyNjA4OGQxZTJhZmZhMTBmZmFjMzBjZjFkMGI5MTViYmI2Yzc5YmZhZmJiMWRiMGU4Y2Q2Njg4MGJmNGJhNTJjNTNiOTQ5ZjZhM2FkYmJhMTgyMWRkMzA0NWM3Y2JcIixcbiAgXCJhODFkOGE5MTc4NjFjYjhhMWZmZDJlOTQ0NTJlMDhmZDZkYzRkMjU3N2JhZDNhYzA4OWY1NjI3OTUyMWIxYzk1Y2FlYmU1N2FjNmMzZDEyNmQ4YWJkNGQ2YTNmMmMzZDhjMjA3YmQzNmZiZjc1YTVlNzk3YzVhOGYxOTkyYzY1MVwiLFxuICBcImI5ZmMzOWQ3ZjEzOGE5NWI4ZjMxNDM2ZTAyYTg3MTFiMzE2NGNkNDVhMjExNjczZjcxMzc0MjliNDVmYWY3N2ExNjA0NjgyYjUxODAzYTk4MzYzOGM0NmE4YjJjMTMyMzdjODdhYjRiNjg1YTFmYTVjNjU3MDBkYzcxMzZjY2ZjXCIsXG4gIFwiMWE1ZjNlYWIwMTI3NmJmN2Q5YjA2YzQyYmU5MGZiNGI4MTA2YjI3OGI0YmJhZjkyMmYzZGE2ODIxYTYzYjRkNjlmYTkyODVmZWMwOTM2ZTBmMDRhMWIyYTI1YTY1MDY0Y2Q1MWMxMTFjNzUwNjNkYmU1OWU0ZGUzMzZiMzVmODVcIixcbiAgXCI1M2JjYjlmZTQwM2E0YTQ1YmVlMmEyYTA0ZGFiZmEzZjIwMThkYjM0OWQwZTUyMDAxNzViZDM0NWFhYTNhM2JkNDVhODhhNmZiNzI0NDkxNGFkMTU2OTYxNzQyZjdiNGVhN2Y4ZmZkODNmY2FlNWIxMTY2YjczYjJhZDk0M2Y3NlwiXG5dXG5cbmludGVyZmFjZSBWZWN0b3Ige1xuICBlbnRyb3B5OiBzdHJpbmdcbiAgbW5lbW9uaWM6IHN0cmluZ1xuICBzZWVkOiBzdHJpbmdcbn1cblxuY29uc3QgdmVjdG9yczogVmVjdG9yW10gPSBbXG4gIHtcbiAgICBlbnRyb3B5OiBcIjAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYm91dFwiLFxuICAgIHNlZWQ6IFwiYzU1MjU3YzM2MGMwN2M3MjAyOWFlYmMxYjUzYzA1ZWQwMzYyYWRhMzhlYWQzZTNlOWVmYTM3MDhlNTM0OTU1MzFmMDlhNjk4NzU5OWQxODI2NGMxZTFjOTJmMmNmMTQxNjMwYzdhM2M0YWI3YzgxYjJmMDAxNjk4ZTc0NjNiMDRcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCI3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZlwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJsZWdhbCB3aW5uZXIgdGhhbmsgeWVhciB3YXZlIHNhdXNhZ2Ugd29ydGggdXNlZnVsIGxlZ2FsIHdpbm5lciB0aGFuayB5ZWxsb3dcIixcbiAgICBzZWVkOiBcIjJlODkwNTgxOWI4NzIzZmUyYzFkMTYxODYwZTVlZTE4MzAzMThkYmY0OWE4M2JkNDUxY2ZiODQ0MGMyOGJkNmZhNDU3ZmUxMjk2MTA2NTU5YTNjODA5MzdhMWMxMDY5YmUzYTNhNWJkMzgxZWU2MjYwZThkOTczOWZjZTFmNjA3XCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODBcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwibGV0dGVyIGFkdmljZSBjYWdlIGFic3VyZCBhbW91bnQgZG9jdG9yIGFjb3VzdGljIGF2b2lkIGxldHRlciBhZHZpY2UgY2FnZSBhYm92ZVwiLFxuICAgIHNlZWQ6IFwiZDcxZGU4NTZmODFhOGFjYzY1ZTZmYzg1MWEzOGQ0ZDdlYzIxNmZkMDc5NmQwYTY4MjdhM2FkNmVkNTUxMWEzMGZhMjgwZjEyZWIyZTQ3ZWQyYWMwM2I1YzQ2MmEwMzU4ZDE4ZDY5ZmU0Zjk4NWVjODE3NzhjMWIzNzBiNjUyYThcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCJmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZlwiLFxuICAgIG1uZW1vbmljOiBcInpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gd3JvbmdcIixcbiAgICBzZWVkOiBcImFjMjc0OTU0ODAyMjUyMjIwNzlkN2JlMTgxNTgzNzUxZTg2ZjU3MTAyN2IwNDk3YjViNWQxMTIxOGUwYThhMTMzMzI1NzI5MTdmMGY4ZTVhNTg5NjIwYzZmMTViMTFjNjFkZWUzMjc2NTFhMTRjMzRlMTgyMzEwNTJlNDhjMDY5XCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhZ2VudFwiLFxuICAgIHNlZWQ6IFwiMDM1ODk1ZjJmNDgxYjFiMGYwMWZjZjhjMjg5Yzc5NDY2MGIyODk5ODFhNzhmODEwNjQ0NzcwN2ZkZDk2NjZjYTA2ZGE1YTlhNTY1MTgxNTk5Yjc5ZjUzYjg0NGQ4YTcxZGQ5ZjQzOWM1MmEzZDdiM2U4YTc5YzkwNmFjODQ1ZmFcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCI3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2ZcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwibGVnYWwgd2lubmVyIHRoYW5rIHllYXIgd2F2ZSBzYXVzYWdlIHdvcnRoIHVzZWZ1bCBsZWdhbCB3aW5uZXIgdGhhbmsgeWVhciB3YXZlIHNhdXNhZ2Ugd29ydGggdXNlZnVsIGxlZ2FsIHdpbGxcIixcbiAgICBzZWVkOiBcImYyYjk0NTA4NzMyYmNiYWNiY2MwMjBmYWVmZWNmYzg5ZmVhZmE2NjQ5YTU0OTFiOGM5NTJjZWRlNDk2YzIxNGEwYzdiM2MzOTJkMTY4NzQ4ZjJkNGE2MTJiYWRhMDc1M2I1MmExYzdhYzUzYzFlOTNhYmQ1YzYzMjBiOWU5NWRkXCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImxldHRlciBhZHZpY2UgY2FnZSBhYnN1cmQgYW1vdW50IGRvY3RvciBhY291c3RpYyBhdm9pZCBsZXR0ZXIgYWR2aWNlIGNhZ2UgYWJzdXJkIGFtb3VudCBkb2N0b3IgYWNvdXN0aWMgYXZvaWQgbGV0dGVyIGFsd2F5c1wiLFxuICAgIHNlZWQ6IFwiMTA3ZDdjMDJhNWFhNmYzOGM1ODA4M2ZmNzRmMDRjNjA3YzJkMmMwZWNjNTU1MDFkYWRkNzJkMDI1Yjc1MWJjMjdmZTkxM2ZmYjc5NmY4NDFjNDliMWQzM2I2MTBjZjBlOTFkM2FhMjM5MDI3ZjVlOTlmZTRjZTllNTA4OGNkNjVcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCJmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwiem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB3aGVuXCIsXG4gICAgc2VlZDogXCIwY2Q2ZTVkODI3YmI2MmViOGZjMWUyNjIyNTQyMjM4MTdmZDA2OGE3NGI1YjQ0OWNjMmY2NjdjM2YxZjk4NWE3NjM3OWI0MzM0OGQ5NTJlMjI2NWI0Y2QxMjkwOTA3NThiM2UzYzJjNDkxMDNiNTA1MWFhYzJlYWViODkwYTUyOFwiXG4gIH0sXG4gIHtcbiAgICBlbnRyb3B5OiBcIjAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDBcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwiYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFydFwiLFxuICAgIHNlZWQ6IFwiYmRhODU0NDZjNjg0MTM3MDcwOTBhNTIwMjJlZGQyNmExYzk0NjIyOTUwMjlmMmU2MGNkN2M0ZjJiYmQzMDk3MTcwYWY3YTRkNzMyNDVjYWZhOWMzY2NhOGQ1NjFhN2MzZGU2ZjVkNGExMGJlOGVkMmE1ZTYwOGQ2OGY5MmZjYzhcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCI3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmN2Y3ZjdmXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImxlZ2FsIHdpbm5lciB0aGFuayB5ZWFyIHdhdmUgc2F1c2FnZSB3b3J0aCB1c2VmdWwgbGVnYWwgd2lubmVyIHRoYW5rIHllYXIgd2F2ZSBzYXVzYWdlIHdvcnRoIHVzZWZ1bCBsZWdhbCB3aW5uZXIgdGhhbmsgeWVhciB3YXZlIHNhdXNhZ2Ugd29ydGggdGl0bGVcIixcbiAgICBzZWVkOiBcImJjMDlmY2ExODA0ZjdlNjlkYTkzYzJmMjAyOGViMjM4YzIyN2YyZTlkZGEzMGNkNjM2OTkyMzI1Nzg0ODBhNDAyMWIxNDZhZDcxN2ZiYjdlNDUxY2U5ZWI4MzVmNDM2MjBiZjVjNTE0ZGIwZjhhZGQ0OWY1ZDEyMTQ0OWQzZTg3XCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MDgwODA4MFwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJsZXR0ZXIgYWR2aWNlIGNhZ2UgYWJzdXJkIGFtb3VudCBkb2N0b3IgYWNvdXN0aWMgYXZvaWQgbGV0dGVyIGFkdmljZSBjYWdlIGFic3VyZCBhbW91bnQgZG9jdG9yIGFjb3VzdGljIGF2b2lkIGxldHRlciBhZHZpY2UgY2FnZSBhYnN1cmQgYW1vdW50IGRvY3RvciBhY291c3RpYyBibGVzc1wiLFxuICAgIHNlZWQ6IFwiYzBjNTE5YmQwZTkxYTJlZDU0MzU3ZDlkMWViZWY2ZjVhZjIxOGExNTM2MjRjZjRmMmRhOTExYTBlZDhmN2EwOWUyZWY2MWFmMGFjYTAwNzA5NmRmNDMwMDIyZjdhMmI2ZmI5MTY2MWE5NTg5MDk3MDY5NzIwZDAxNWU0ZTk4MmZcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCJmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmZmXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcInpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gdm90ZVwiLFxuICAgIHNlZWQ6IFwiZGQ0OGMxMDQ2OThjMzBjZmUyYjYxNDIxMDMyNDg2MjJmYjdiYjBmZjY5MmVlYmIwMDA4OWIzMmQyMjQ4NGUxNjEzOTEyZjBhNWI2OTQ0MDdiZTg5OWZmZDMxZWQzOTkyYzQ1NmNkZjYwZjVkNDU2NGI4YmEzZjA1YTY5ODkwYWRcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCI3N2MyYjAwNzE2Y2VjNzIxMzgzOTE1OWU0MDRkYjUwZFwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJqZWxseSBiZXR0ZXIgYWNoaWV2ZSBjb2xsZWN0IHVuYXdhcmUgbW91bnRhaW4gdGhvdWdodCBjYXJnbyBveHlnZW4gYWN0IGhvb2QgYnJpZGdlXCIsXG4gICAgc2VlZDogXCJiNWI2ZDAxMjdkYjFhOWQyMjI2YWYwYzMzNDYwMzFkNzdhZjMxZTkxOGRiYTY0Mjg3YTFiNDRiOGViZjYzY2RkNTI2NzZmNjcyYTI5MGFhZTUwMjQ3MmNmMmQ2MDJjMDUxZjNlNmYxODA1NWU4NGU0YzQzODk3ZmM0ZTUxYTZmZlwiXG4gIH0sXG4gIHtcbiAgICBlbnRyb3B5OiBcImI2M2E5YzU5YTZlNjQxZjI4OGViYzEwMzAxN2YxZGE5ZjgyOTBiM2RhNmJkZWY3YlwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJyZW5ldyBzdGF5IGJpb2xvZ3kgZXZpZGVuY2UgZ29hdCB3ZWxjb21lIGNhc3VhbCBqb2luIGFkYXB0IGFybW9yIHNodWZmbGUgZmF1bHQgbGl0dGxlIG1hY2hpbmUgd2FsayBzdHVtYmxlIHVyZ2Ugc3dhcFwiLFxuICAgIHNlZWQ6IFwiOTI0OGQ4M2UwNmY0Y2Q5OGRlYmY1YjZmMDEwNTQyNzYwZGY5MjVjZTQ2Y2YzOGExYmRiNGU0ZGU3ZDIxZjVjMzkzNjY5NDFjNjllMWJkYmYyOTY2ZTBmNmU2ZGJlY2U4OThhMGUyZjBhNGMyYjNlNjQwOTUzZGZlOGI3YmJkYzVcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCIzZTE0MTYwOWI5NzkzM2I2NmEwNjBkY2RkYzcxZmFkMWQ5MTY3N2RiODcyMDMxZTg1ZjRjMDE1YzVlN2U4OTgyXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImRpZ25pdHkgcGFzcyBsaXN0IGluZGljYXRlIG5hc3R5IHN3YW1wIHBvb2wgc2NyaXB0IHNvY2NlciB0b2UgbGVhZiBwaG90byBtdWx0aXBseSBkZXNrIGhvc3QgdG9tYXRvIGNyYWRsZSBkcmlsbCBzcHJlYWQgYWN0b3Igc2hpbmUgZGlzbWlzcyBjaGFtcGlvbiBleG90aWNcIixcbiAgICBzZWVkOiBcImZmN2YzMTg0ZGY4Njk2ZDhiZWY5NGI2YzAzMTE0ZGJlZTBlZjg5ZmY5Mzg3MTIzMDFkMjdlZDgzMzZjYTg5ZWY5NjM1ZGEyMGFmMDdkNDE3NWYyYmY1ZjNkZTEzMGYzOWM5ZDllOGRkMDQ3MjQ4OWMxOWIxYTAyMGE5NDBkYTY3XCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiMDQ2MGVmNDc1ODU2MDRjNTY2MDYxOGRiMmU2YTdlN2ZcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwiYWZmb3JkIGFsdGVyIHNwaWtlIHJhZGFyIGdhdGUgZ2xhbmNlIG9iamVjdCBzZWVrIHN3YW1wIGluZmFudCBwYW5lbCB5ZWxsb3dcIixcbiAgICBzZWVkOiBcIjY1ZjkzYTlmMzZiNmM4NWNiZTYzNGZmYzFmOTlmMmI4MmNiYjEwYjMxZWRjN2YwODdiNGY2Y2I5ZTk3NmU5ZmFmNzZmZjQxZjhmMjdjOTlhZmRmMzhmN2EzMDNiYTExMzZlZTQ4YTRjMWU3ZmNkM2RiYTdhYTg3NjExM2EzNmU0XCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiNzJmNjBlYmFjNWRkOGFkZDhkMmEyNWE3OTcxMDJjM2NlMjFiYzAyOWMyMDAwNzZmXCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImluZGljYXRlIHJhY2UgcHVzaCBtZXJyeSBzdWZmZXIgaHVtYW4gY3J1aXNlIGR3YXJmIHBvbGUgcmV2aWV3IGFyY2gga2VlcCBjYW52YXMgdGhlbWUgcG9lbSBkaXZvcmNlIGFsdGVyIGxlZnRcIixcbiAgICBzZWVkOiBcIjNiYmY5ZGFhMGRmYWQ4MjI5Nzg2YWNlNWRkYjRlMDBmYTk4YTA0NGFlNGM0OTc1ZmZkNWUwOTRkYmE5ZTBiYjI4OTM0OWRiZTIwOTE3NjFmMzBmMzgyZDRlMzVjNGE2NzBlZThhYjUwNzU4ZDJjNTU4ODFiZTY5ZTMyNzExN2JhXCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiMmM4NWVmYzdmMjRlZTQ1NzNkMmI4MWE2ZWM2NmNlZTIwOWIyZGNiZDA5ZDhlZGRjNTFlMDIxNWIwYjY4ZTQxNlwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJjbHV0Y2ggY29udHJvbCB2ZWhpY2xlIHRvbmlnaHQgdW51c3VhbCBjbG9nIHZpc2EgaWNlIHBsdW5nZSBnbGltcHNlIHJlY2lwZSBzZXJpZXMgb3BlbiBob3VyIHZpbnRhZ2UgZGVwb3NpdCB1bml2ZXJzZSB0aXAgam9iIGRyZXNzIHJhZGFyIHJlZnVzZSBtb3Rpb24gdGFzdGVcIixcbiAgICBzZWVkOiBcImZlOTA4Zjk2ZjQ2NjY4YjJkNWIzN2Q4MmY1NThjNzdlZDBkNjlkZDBlN2UwNDNhNWIwNTExYzQ4YzJmMTA2NDY5NGE5NTZmODYzNjBjOTNkZDA0MDUyYTg4OTk0OTdjZTllOTg1ZWJlMGM4YzUyYjk1NWU2YWU4NmQ0ZmY0NDQ5XCJcbiAgfSxcbiAge1xuICAgIGVudHJvcHk6IFwiZWFlYmFiYjIzODMzNTFmZDMxZDcwMzg0MGIzMmU5ZTJcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwidHVydGxlIGZyb250IHVuY2xlIGlkZWEgY3J1c2ggd3JpdGUgc2hydWcgdGhlcmUgbG90dGVyeSBmbG93ZXIgcmlzayBzaGVsbFwiLFxuICAgIHNlZWQ6IFwiYmRmYjc2YTA3NTlmMzAxYjBiODk5YTFlMzk4NTIyN2U1M2IzZjUxZTY3ZTNmMmE2NTM2M2NhZWRmM2UzMmZkZTQyYTY2YzQwNGYxOGQ3YjA1ODE4Yzk1ZWYzY2ExZTUxNDY2NDY4NTZjNDYxYzA3MzE2OTQ2NzUxMTY4MDg3NmNcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCI3YWM0NWNmZTc3MjJlZTZjN2JhODRmYmMyZDViZDYxYjQ1Y2IyZmU1ZWI2NWFhNzhcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwia2lzcyBjYXJyeSBkaXNwbGF5IHVudXN1YWwgY29uZmlybSBjdXJ0YWluIHVwZ3JhZGUgYW50aXF1ZSByb3RhdGUgaGVsbG8gdm9pZCBjdXN0b20gZnJlcXVlbnQgb2JleSBudXQgaG9sZSBwcmljZSBzZWdtZW50XCIsXG4gICAgc2VlZDogXCJlZDU2ZmY2YzgzM2MwNzk4MmViNzExOWE4ZjQ4ZmQzNjNjNGE5YjE2MDFjZDJkZTczNmIwMTA0NWM1ZWI4YWI0ZjU3YjA3OTQwMzQ4NWQxYzQ5MjRmMDc5MGRjMTBhOTcxNzYzMzM3Y2I5ZjljNjIyMjZmNjRmZmYyNjM5N2M3OVwiXG4gIH0sXG4gIHtcbiAgICBlbnRyb3B5OiBcIjRmYTFhOGJjM2U2ZDgwZWUxMzE2MDUwZTg2MmMxODEyMDMxNDkzMjEyYjdlYzNmM2JiMWIwOGYxNjhjYWJlZWZcIixcbiAgICBtbmVtb25pYzpcbiAgICAgIFwiZXhpbGUgYXNrIGNvbmdyZXNzIGxhbXAgc3VibWl0IGphY2tldCBlcmEgc2NoZW1lIGF0dGVuZCBjb3VzaW4gYWxjb2hvbCBjYXRjaCBjb3Vyc2UgZW5kIGx1Y2t5IGh1cnQgc2VudGVuY2Ugb3ZlbiBzaG9ydCBiYWxsIGJpcmQgZ3JhYiB3aW5nIHRvcFwiLFxuICAgIHNlZWQ6IFwiMDk1ZWU2ZjgxN2I0YzJjYjMwYTVhNzk3MzYwYTgxYTQwYWIwZjlhNGUyNWVjZDY3MmEzZjU4YTBiNWJhMDY4N2MwOTZhNmIxNGQyYzBkZWIzYmRlZmNlNGY2MWQwMWFlMDc0MTdkNTAyNDI5MzUyZTI3Njk1MTYzZjc0NDdhOGNcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCIxOGFiMTlhOWY1NGE5Mjc0ZjAzZTUyMDlhMmFjOGE5MVwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJib2FyZCBmbGVlIGhlYXZ5IHR1bm5lbCBwb3dkZXIgZGVuaWFsIHNjaWVuY2Ugc2tpIGFuc3dlciBiZXRyYXkgY2FyZ28gY2F0XCIsXG4gICAgc2VlZDogXCI2ZWZmMWJiMjE1NjI5MTg1MDljNzNjYjk5MDI2MGRiMDdjMGNlMzRmZjBlM2NjNGE4Y2IzMjc2MTI5ZmJjYjMwMGJkZGZlMDA1ODMxMzUwZWZkNjMzOTA5ZjQ3NmM0NWM4ODI1MzI3NmQ5ZmQwZGY2ZWY0ODYwOWU4YmI3ZGNhOFwiXG4gIH0sXG4gIHtcbiAgICBlbnRyb3B5OiBcIjE4YTJlMWQ4MWI4ZWNmYjJhMzMzYWRjYjBjMTdhNWI5ZWI3NmNjNWQwNWRiOTFhNFwiLFxuICAgIG1uZW1vbmljOlxuICAgICAgXCJib2FyZCBibGFkZSBpbnZpdGUgZGFtYWdlIHVuZG8gc3VuIG1pbWljIGludGVyZXN0IHNsYW0gZ2F6ZSB0cnVseSBpbmhlcml0IHJlc2lzdCBncmVhdCBpbmplY3Qgcm9ja2V0IG11c2V1bSBjaGllZlwiLFxuICAgIHNlZWQ6IFwiZjg0NTIxYzc3N2ExM2I2MTU2NDIzNGJmOGY4YjYyYjNhZmNlMjdmYzQwNjJiNTFiYjVlNjJiZGZlY2IyMzg2NGVlNmVjZjA3YzFkNWE5N2MwODM0MzA3YzVjODUyZDhjZWI4OGU3Yzk3OTIzYzBhM2I0OTZiZWRkNGU1Zjg4YTlcIlxuICB9LFxuICB7XG4gICAgZW50cm9weTogXCIxNWRhODcyYzk1YTEzZGQ3MzhmYmY1MGU0Mjc1ODNhZDYxZjE4ZmQ5OWY2MjhjNDE3YTYxY2Y4MzQzYzkwNDE5XCIsXG4gICAgbW5lbW9uaWM6XG4gICAgICBcImJleW9uZCBzdGFnZSBzbGVlcCBjbGlwIGJlY2F1c2UgdHdpc3QgdG9rZW4gbGVhZiBhdG9tIGJlYXV0eSBnZW5pdXMgZm9vZCBidXNpbmVzcyBzaWRlIGdyaWQgdW5hYmxlIG1pZGRsZSBhcm1lZCBvYnNlcnZlIHBhaXIgY3JvdWNoIHRvbmlnaHQgYXdheSBjb2NvbnV0XCIsXG4gICAgc2VlZDogXCJiMTU1MDllYWEyZDA5ZDNlZmQzZTAwNmVmNDIxNTFiMzAzNjdkYzZlM2FhNWU0NGNhYmEzZmU0ZDNlMzUyZTY1MTAxZmJkYjg2YTk2Nzc2YjkxOTQ2ZmYwNmY4ZWFjNTk0ZGM2ZWUxZDNlODJhNDJkZmUxYjQwZmVmNmJjYzNmZFwiXG4gIH1cbl1cblxuaW50ZXJmYWNlIEJhZE1uZW1vbmljIHtcbiAgbW5lbW9uaWM6IHN0cmluZ1xufVxuXG5jb25zdCBiYWRNbmVtb25pY3M6IEJhZE1uZW1vbmljW10gPSBbXG4gIHtcbiAgICBtbmVtb25pYzpcbiAgICAgIFwiYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uXCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJsZWdhbCB3aW5uZXIgdGhhbmsgeWVhciB3YXZlIHNhdXNhZ2Ugd29ydGggdXNlZnVsIGxlZ2FsIHdpbm5lciB0aGFuayB5ZWxsb3cgeWVsbG93XCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJsZXR0ZXIgYWR2aWNlIGNhZ2UgYWJzdXJkIGFtb3VudCBkb2N0b3IgYWNvdXN0aWMgYXZvaWQgbGV0dGVyIGFkdmljZSBjYWdlZCBhYm92ZVwiXG4gIH0sXG4gIHsgbW5lbW9uaWM6IFwiem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbywgd3JvbmdcIiB9LFxuICB7XG4gICAgbW5lbW9uaWM6XG4gICAgICBcImFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvblwiXG4gIH0sXG4gIHtcbiAgICBtbmVtb25pYzpcbiAgICAgIFwibGVnYWwgd2lubmVyIHRoYW5rIHllYXIgd2F2ZSBzYXVzYWdlIHdvcnRoIHVzZWZ1bCBsZWdhbCB3aW5uZXIgdGhhbmsgeWVhciB3YXZlIHNhdXNhZ2Ugd29ydGggdXNlZnVsIGxlZ2FsIHdpbGwgd2lsbCB3aWxsXCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJsZXR0ZXIgYWR2aWNlIGNhZ2UgYWJzdXJkIGFtb3VudCBkb2N0b3IgYWNvdXN0aWMgYXZvaWQgbGV0dGVyIGFkdmljZSBjYWdlIGFic3VyZCBhbW91bnQgZG9jdG9yIGFjb3VzdGljIGF2b2lkIGxldHRlciBhbHdheXMuXCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJ6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHdoeVwiXG4gIH0sXG4gIHtcbiAgICBtbmVtb25pYzpcbiAgICAgIFwiYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFydCBhcnRcIlxuICB9LFxuICB7XG4gICAgbW5lbW9uaWM6XG4gICAgICBcImxlZ2FsIHdpbm5lciB0aGFuayB5ZWFyIHdhdmUgc2F1c2FnZSB3b3J0aCB1c2VmdWwgbGVnYWwgd2lubmVyIHRoYW5rcyB5ZWFyIHdhdmUgd29ydGggdXNlZnVsIGxlZ2FsIHdpbm5lciB0aGFuayB5ZWFyIHdhdmUgc2F1c2FnZSB3b3J0aCB0aXRsZVwiXG4gIH0sXG4gIHtcbiAgICBtbmVtb25pYzpcbiAgICAgIFwibGV0dGVyIGFkdmljZSBjYWdlIGFic3VyZCBhbW91bnQgZG9jdG9yIGFjb3VzdGljIGF2b2lkIGxldHRlcnMgYWR2aWNlIGNhZ2UgYWJzdXJkIGFtb3VudCBkb2N0b3IgYWNvdXN0aWMgYXZvaWQgbGV0dGVyIGFkdmljZSBjYWdlIGFic3VyZCBhbW91bnQgZG9jdG9yIGFjb3VzdGljIGJsZXNzXCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJ6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHpvbyB6b28gem9vIHZvdGVkXCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJqZWxsbyBiZXR0ZXIgYWNoaWV2ZSBjb2xsZWN0IHVuYXdhcmUgbW91bnRhaW4gdGhvdWdodCBjYXJnbyBveHlnZW4gYWN0IGhvb2QgYnJpZGdlXCJcbiAgfSxcbiAge1xuICAgIG1uZW1vbmljOlxuICAgICAgXCJyZW5ldywgc3RheSwgYmlvbG9neSwgZXZpZGVuY2UsIGdvYXQsIHdlbGNvbWUsIGNhc3VhbCwgam9pbiwgYWRhcHQsIGFybW9yLCBzaHVmZmxlLCBmYXVsdCwgbGl0dGxlLCBtYWNoaW5lLCB3YWxrLCBzdHVtYmxlLCB1cmdlLCBzd2FwXCJcbiAgfSxcbiAgeyBtbmVtb25pYzogXCJkaWduaXR5IHBhc3MgbGlzdCBpbmRpY2F0ZSBuYXN0eVwiIH0sXG5cbiAgLy8gRnJvbSBpc3N1ZSAzMlxuICB7XG4gICAgbW5lbW9uaWM6XG4gICAgICBcImFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBhYmFuZG9uIGFiYW5kb24gYWJhbmRvbiBsZXR0ZXJcIlxuICB9XG5dXG5cbmNvbnN0IG1hbGZvcm1lZE1uZW1vbmljczogc3RyaW5nW10gPSBbXG4gIFwiYSBhIGEgYSBhIGEgYSBhIGEgYSBhIGEgYSBhIGEgYSBhIGEgYSBhIGEgYSBhIGEgYVwiLCAvLyBUb28gbWFueSB3b3Jkc1xuICBcImFcIiwgLy8gVG9vIGZld1xuICBcImEgYSBhIGEgYSBhIGEgYSBhIGEgYSBhIGEgYVwiIC8vIE5vdCBtdWx0aXBsZSBvZiAzXG5dXG5cbmRlc2NyaWJlKFwiTW5lbW9uaWNcIiwgKCkgPT4ge1xuICB0ZXN0KFwidmVjdG9yc1wiLCAoKTogdm9pZCA9PiB7XG4gICAgdmVjdG9ycy5mb3JFYWNoKGFzeW5jICh2ZWN0b3I6IFZlY3RvciwgaW5kZXg6IG51bWJlcik6IFByb21pc2U8dm9pZD4gPT4ge1xuICAgICAgY29uc3Qgd29yZGxpc3QgPSBtbmVtb25pYy5nZXRXb3JkbGlzdHMoXCJlbmdsaXNoXCIpIGFzIHN0cmluZ1tdXG4gICAgICBjb25zdCBlbnRyb3B5VG9NbmVtb25pYzogc3RyaW5nID0gbW5lbW9uaWMuZW50cm9weVRvTW5lbW9uaWMoXG4gICAgICAgIHZlY3Rvci5lbnRyb3B5LFxuICAgICAgICB3b3JkbGlzdFxuICAgICAgKVxuICAgICAgZXhwZWN0KHZlY3Rvci5tbmVtb25pYykudG9CZShlbnRyb3B5VG9NbmVtb25pYylcbiAgICAgIGNvbnN0IG1uZW1vbmljVG9FbnRyb3B5OiBzdHJpbmcgPSBtbmVtb25pYy5tbmVtb25pY1RvRW50cm9weShcbiAgICAgICAgdmVjdG9yLm1uZW1vbmljLFxuICAgICAgICB3b3JkbGlzdFxuICAgICAgKVxuICAgICAgZXhwZWN0KG1uZW1vbmljVG9FbnRyb3B5KS50b0JlKHZlY3Rvci5lbnRyb3B5KVxuICAgICAgY29uc3QgcGFzc3dvcmQ6IHN0cmluZyA9IFwiVFJFWk9SXCJcbiAgICAgIGNvbnN0IG1uZW1vbmljVG9TZWVkOiBCdWZmZXIgPSBhd2FpdCBtbmVtb25pYy5tbmVtb25pY1RvU2VlZChcbiAgICAgICAgdmVjdG9yLm1uZW1vbmljLFxuICAgICAgICBwYXNzd29yZFxuICAgICAgKVxuICAgICAgZXhwZWN0KG1uZW1vbmljVG9TZWVkLnRvU3RyaW5nKFwiaGV4XCIpKS50b0JlKHZlY3Rvci5zZWVkKVxuICAgIH0pXG4gIH0pXG5cbiAgdGVzdChcImJhZE1uZW1vbmljc1wiLCAoKTogdm9pZCA9PiB7XG4gICAgY29uc3Qgd29yZGxpc3QgPSBtbmVtb25pYy5nZXRXb3JkbGlzdHMoXCJlbmdsaXNoXCIpIGFzIHN0cmluZ1tdXG4gICAgYmFkTW5lbW9uaWNzLmZvckVhY2goKGJhZE1uZW1vbmljOiBCYWRNbmVtb25pYywgaW5kZXg6IG51bWJlcik6IHZvaWQgPT4ge1xuICAgICAgY29uc3QgdmFsaWRhdGVNbmVtb25pYzogc3RyaW5nID0gbW5lbW9uaWMudmFsaWRhdGVNbmVtb25pYyhcbiAgICAgICAgYmFkTW5lbW9uaWMubW5lbW9uaWMsXG4gICAgICAgIHdvcmRsaXN0XG4gICAgICApXG4gICAgICBleHBlY3QodmFsaWRhdGVNbmVtb25pYykudG9CZUZhbHN5KClcbiAgICB9KVxuICB9KVxuXG4gIHRlc3QoXCJtYWxmb3JtZWRNbmVtb25pY3NcIiwgKCk6IHZvaWQgPT4ge1xuICAgIGNvbnN0IHdvcmRsaXN0ID0gbW5lbW9uaWMuZ2V0V29yZGxpc3RzKFwiZW5nbGlzaFwiKSBhcyBzdHJpbmdbXVxuICAgIG1hbGZvcm1lZE1uZW1vbmljcy5mb3JFYWNoKFxuICAgICAgKG1hbGZvcm1lZE1uZW1vbmljOiBzdHJpbmcsIGluZGV4OiBudW1iZXIpOiB2b2lkID0+IHtcbiAgICAgICAgY29uc3QgdmFsaWRhdGVNbmVtb25pYzogc3RyaW5nID0gbW5lbW9uaWMudmFsaWRhdGVNbmVtb25pYyhcbiAgICAgICAgICBtYWxmb3JtZWRNbmVtb25pYyxcbiAgICAgICAgICB3b3JkbGlzdFxuICAgICAgICApXG4gICAgICAgIGV4cGVjdCh2YWxpZGF0ZU1uZW1vbmljKS50b0JlRmFsc3koKVxuICAgICAgfVxuICAgIClcbiAgfSlcblxuICB0ZXN0KFwiZW50cm9weVRvTW5lbW9uaWNcIiwgKCk6IHZvaWQgPT4ge1xuICAgIGxhbmdzLmZvckVhY2goKGxhbmc6IHN0cmluZywgaW5kZXg6IG51bWJlcik6IHZvaWQgPT4ge1xuICAgICAgY29uc3Qgd29yZGxpc3QgPSBtbmVtb25pYy5nZXRXb3JkbGlzdHMobGFuZykgYXMgc3RyaW5nW11cbiAgICAgIGNvbnN0IGVudHJvcHlUb01uZW1vbmljOiBzdHJpbmcgPSBtbmVtb25pYy5lbnRyb3B5VG9NbmVtb25pYyhcbiAgICAgICAgZW50cm9weSxcbiAgICAgICAgd29yZGxpc3RcbiAgICAgIClcbiAgICAgIGV4cGVjdChtbmVtbm9uaWNzW2luZGV4XSkudG9CZShlbnRyb3B5VG9NbmVtb25pYylcbiAgICB9KVxuICB9KVxuXG4gIHRlc3QoXCJnZW5lcmF0ZU1uZW1vbmljXCIsICgpOiB2b2lkID0+IHtcbiAgICBjb25zdCBzdHJlbmd0aDogbnVtYmVyID0gMjU2XG4gICAgbGFuZ3MuZm9yRWFjaCgobGFuZzogc3RyaW5nKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB3b3JkbGlzdCA9IG1uZW1vbmljLmdldFdvcmRsaXN0cyhsYW5nKSBhcyBzdHJpbmdbXVxuICAgICAgY29uc3QgbTogc3RyaW5nID0gbW5lbW9uaWMuZ2VuZXJhdGVNbmVtb25pYyhcbiAgICAgICAgc3RyZW5ndGgsXG4gICAgICAgIHJhbmRvbUJ5dGVzLFxuICAgICAgICB3b3JkbGlzdFxuICAgICAgKVxuICAgICAgZXhwZWN0KHR5cGVvZiBtID09PSBcInN0cmluZ1wiKS50b0JlVHJ1dGh5KClcbiAgICB9KVxuICB9KVxuXG4gIHRlc3QoXCJ0ZXN0IG1uZW1vbmljIGxlbmd0aHNcIiwgKCk6IHZvaWQgPT4ge1xuICAgIGNvbnN0IHdvcmRsaXN0ID0gbW5lbW9uaWMuZ2V0V29yZGxpc3RzKFwiZW5nbGlzaFwiKSBhcyBzdHJpbmdbXVxuICAgIGxldCBzdHJlbmd0aDogbnVtYmVyID0gMTI4XG4gICAgbGV0IG1uZW1ubmljOiBzdHJpbmcgPSBtbmVtb25pYy5nZW5lcmF0ZU1uZW1vbmljKFxuICAgICAgc3RyZW5ndGgsXG4gICAgICByYW5kb21CeXRlcyxcbiAgICAgIHdvcmRsaXN0XG4gICAgKVxuICAgIGV4cGVjdChtbmVtbm5pYy5zcGxpdChcIiBcIikubGVuZ3RoKS50b0JlKDEyKVxuICAgIHN0cmVuZ3RoID0gMTYwXG4gICAgbW5lbW5uaWMgPSBtbmVtb25pYy5nZW5lcmF0ZU1uZW1vbmljKHN0cmVuZ3RoLCByYW5kb21CeXRlcywgd29yZGxpc3QpXG4gICAgZXhwZWN0KG1uZW1ubmljLnNwbGl0KFwiIFwiKS5sZW5ndGgpLnRvQmUoMTUpXG4gICAgc3RyZW5ndGggPSAxOTJcbiAgICBtbmVtbm5pYyA9IG1uZW1vbmljLmdlbmVyYXRlTW5lbW9uaWMoc3RyZW5ndGgsIHJhbmRvbUJ5dGVzLCB3b3JkbGlzdClcbiAgICBleHBlY3QobW5lbW5uaWMuc3BsaXQoXCIgXCIpLmxlbmd0aCkudG9CZSgxOClcbiAgICBzdHJlbmd0aCA9IDIyNFxuICAgIG1uZW1ubmljID0gbW5lbW9uaWMuZ2VuZXJhdGVNbmVtb25pYyhzdHJlbmd0aCwgcmFuZG9tQnl0ZXMsIHdvcmRsaXN0KVxuICAgIGV4cGVjdChtbmVtbm5pYy5zcGxpdChcIiBcIikubGVuZ3RoKS50b0JlKDIxKVxuICAgIHN0cmVuZ3RoID0gMjU2XG4gICAgbW5lbW5uaWMgPSBtbmVtb25pYy5nZW5lcmF0ZU1uZW1vbmljKHN0cmVuZ3RoLCByYW5kb21CeXRlcywgd29yZGxpc3QpXG4gICAgZXhwZWN0KG1uZW1ubmljLnNwbGl0KFwiIFwiKS5sZW5ndGgpLnRvQmUoMjQpXG4gIH0pXG5cbiAgdGVzdChcImdldFdvcmRsaXN0c1wiLCAoKTogdm9pZCA9PiB7XG4gICAgbGFuZ3MuZm9yRWFjaCgobGFuZzogc3RyaW5nKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCB3b3JkbGlzdCA9IG1uZW1vbmljLmdldFdvcmRsaXN0cyhsYW5nKSBhcyBzdHJpbmdbXVxuICAgICAgZXhwZWN0KHR5cGVvZiB3b3JkbGlzdCA9PT0gXCJvYmplY3RcIikudG9CZVRydXRoeSgpXG4gICAgfSlcbiAgfSlcblxuICB0ZXN0KFwibW5lbW9uaWNUb0VudHJvcHlcIiwgKCk6IHZvaWQgPT4ge1xuICAgIG1uZW1ub25pY3MuZm9yRWFjaCgobW5lbW5uaWM6IHN0cmluZywgaW5kZXg6IG51bWJlcik6IHZvaWQgPT4ge1xuICAgICAgY29uc3Qgd29yZGxpc3QgPSBtbmVtb25pYy5nZXRXb3JkbGlzdHMobGFuZ3NbaW5kZXhdKSBhcyBzdHJpbmdbXVxuICAgICAgY29uc3QgbW5lbW9uaWNUb0VudHJvcHk6IHN0cmluZyA9IG1uZW1vbmljLm1uZW1vbmljVG9FbnRyb3B5KFxuICAgICAgICBtbmVtbm5pYyxcbiAgICAgICAgd29yZGxpc3RcbiAgICAgIClcbiAgICAgIGV4cGVjdChtbmVtb25pY1RvRW50cm9weSkudG9CZShlbnRyb3B5KVxuICAgIH0pXG4gIH0pXG5cbiAgdGVzdChcIm1uZW1vbmljVG9TZWVkXCIsIGFzeW5jICgpOiBQcm9taXNlPHZvaWQ+ID0+IHtcbiAgICBtbmVtbm9uaWNzLmZvckVhY2goYXN5bmMgKG1uZW1ubmljOiBzdHJpbmcpOiBQcm9taXNlPGFueT4gPT4ge1xuICAgICAgY29uc3QgcGFzc3dvcmQ6IHN0cmluZyA9IFwicGFzc3dvcmRcIlxuICAgICAgY29uc3QgbW5lbW9uaWNUb1NlZWQ6IEJ1ZmZlciA9IGF3YWl0IG1uZW1vbmljLm1uZW1vbmljVG9TZWVkKFxuICAgICAgICBtbmVtbm5pYyxcbiAgICAgICAgcGFzc3dvcmRcbiAgICAgIClcbiAgICAgIGV4cGVjdCh0eXBlb2YgbW5lbW9uaWNUb1NlZWQgPT09IFwib2JqZWN0XCIpLnRvQmVUcnV0aHkoKVxuICAgIH0pXG4gIH0pXG5cbiAgdGVzdChcIm1uZW1vbmljVG9TZWVkU3luY1wiLCAoKTogdm9pZCA9PiB7XG4gICAgbW5lbW5vbmljcy5mb3JFYWNoKChtbmVtbm5pYzogc3RyaW5nLCBpbmRleDogbnVtYmVyKTogdm9pZCA9PiB7XG4gICAgICBjb25zdCBwYXNzd29yZDogc3RyaW5nID0gXCJwYXNzd29yZFwiXG4gICAgICBjb25zdCBtbmVtb25pY1RvU2VlZFN5bmM6IEJ1ZmZlciA9IG1uZW1vbmljLm1uZW1vbmljVG9TZWVkU3luYyhcbiAgICAgICAgbW5lbW5uaWMsXG4gICAgICAgIHBhc3N3b3JkXG4gICAgICApXG4gICAgICBleHBlY3QobW5lbW9uaWNUb1NlZWRTeW5jLnRvU3RyaW5nKFwiaGV4XCIpKS50b0JlKHNlZWRzW2luZGV4XSlcbiAgICB9KVxuICB9KVxuXG4gIHRlc3QoXCJ2YWxpZGF0ZU1uZW1vbmljXCIsICgpOiB2b2lkID0+IHtcbiAgICBtbmVtbm9uaWNzLmZvckVhY2goKG1uZW1ubmljOiBzdHJpbmcsIGluZGV4OiBudW1iZXIpOiB2b2lkID0+IHtcbiAgICAgIGNvbnN0IHdvcmRsaXN0ID0gbW5lbW9uaWMuZ2V0V29yZGxpc3RzKGxhbmdzW2luZGV4XSkgYXMgc3RyaW5nW11cbiAgICAgIGNvbnN0IHZhbGlkYXRlTW5lbW9uaWM6IHN0cmluZyA9IG1uZW1vbmljLnZhbGlkYXRlTW5lbW9uaWMoXG4gICAgICAgIG1uZW1ubmljLFxuICAgICAgICB3b3JkbGlzdFxuICAgICAgKVxuICAgICAgZXhwZWN0KHZhbGlkYXRlTW5lbW9uaWMpLnRvQmVUcnV0aHkoKVxuICAgIH0pXG4gIH0pXG5cbiAgdGVzdChcInNldERlZmF1bHRXb3JkbGlzdFwiLCAoKTogdm9pZCA9PiB7XG4gICAgbGFuZ3MuZm9yRWFjaCgobGFuZzogc3RyaW5nLCBpbmRleDogbnVtYmVyKTogdm9pZCA9PiB7XG4gICAgICBtbmVtb25pYy5zZXREZWZhdWx0V29yZGxpc3QobGFuZylcbiAgICAgIGNvbnN0IGdldERlZmF1bHRXb3JkbGlzdDogc3RyaW5nID0gbW5lbW9uaWMuZ2V0RGVmYXVsdFdvcmRsaXN0KClcbiAgICAgIGNvbnN0IHdvcmRsaXN0ID0gbW5lbW9uaWMuZ2V0V29yZGxpc3RzKGxhbmcpIGFzIHN0cmluZ1tdXG4gICAgICBjb25zdCBtOiBzdHJpbmcgPSBtbmVtb25pYy5nZW5lcmF0ZU1uZW1vbmljKDI1NiwgcmFuZG9tQnl0ZXMsIHdvcmRsaXN0KVxuICAgICAgZXhwZWN0KGdldERlZmF1bHRXb3JkbGlzdCkudG9CZShsYW5nKVxuICAgICAgZXhwZWN0KHR5cGVvZiB3b3JkbGlzdCA9PT0gXCJvYmplY3RcIikudG9CZVRydXRoeSgpXG4gICAgICBleHBlY3QodHlwZW9mIG0gPT09IFwic3RyaW5nXCIpLnRvQmVUcnV0aHkoKVxuICAgIH0pXG4gIH0pXG59KVxuIl19