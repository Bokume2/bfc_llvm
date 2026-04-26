# Brainf\*ck Compiler(Processer) based on LLVM
「15日間LLVM生活」の学習用リポジトリです。  

||技術構成|
|:---:|:----|
|フロントエンド|Go|
|バックエンド|LLVM|

## Usage
※学習用プロダクトのため、予告無く破壊的変更が加わる可能性が高いことにご留意ください。  

1. LLVMをインストールします。  
   動作確認済みのバージョン: LLVM 20 from [apt.llvm.org](https://apt.llvm.org) or Ubuntu repository
   ```bash
   # 例: APTからインストールする場合
   # Debian trixie等の場合、apt.llvm.orgのサードパーティリポジトリの登録が必要
   apt install llvm-20-dev
   ```
1. Goをインストールします。  
1. リポジトリをcloneし、./cmd/bfcをビルドします。  
   ```bash
   git clone https://github.com/Bokume2/bfc_llvm.git
   cd bfc_llvm
   go build ./cmd/bfc
   ```
1. Brainf\*ckのソースコードをテキストファイルに記述し、ビルドしたBrainf\*ckコンパイラに渡します。コンパイルされたLLVM IRは標準出力に出力されるため、適宜リダイレクト等してください。  
   ```bash
   ./bfc hoge.bf > hoge.ll
   ```
   オプションとして、LLVMの最適化のレベルを`-O0`～`-O3`で指定することが可能です。デフォルトは`-O3`で、レベルごとの最適化の内容はLLVMの各種コマンドと同じです。  
   ```bash
   ./bfc -O1 hoge.bf > hoge.ll
   ```
1. 生成したLLVM IRをClangで実行可能ファイルにコンパイルします。ターゲットマシンの上書き警告(`Woverride-module`)は基本的に無視して構わないでしょう。  
   ```bash
   clang --Wno-override-module -o hoge hoge.ll
   ```
   あるいは、`llc`でアセンブリにコンパイルしてから、任意のCコンパイラで実行可能ファイルにコンパイルします。  
   ```bash
   llc -o hoge.s hoge.ll
   cc -o hoge hoge.s
   ```
1. 実行可能ファイルを実行します。  
   ```bash
   ./hoge
   ```

## License
ソースコードはMIT Licenseの下に配布されています。  
詳しくは[ライセンス表示](./LICENSE)や[https://www.tldrlegal.com/license/mit-license](https://www.tldrlegal.com/license/mit-license)などを参照>して下さい。  
なお、ライセンスは予告無く変更される可能性があり、変更前のライセンスは、そのライセンス表示がこのリポジトリのmainブランチの最新コミットに適用されていた期間に公開されたファイルにのみ適用されます。  
