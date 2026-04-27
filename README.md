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
1. 実行可能ファイルを、以下のいずれかの方法で生成します。  
   - Cコンパイラ等を手動で実行する方法
     1. Brainf\*ckのソースコードをテキストファイルに記述し、ビルドしたBrainf\*ckコンパイラに渡します。`-o`オプションを省略するとLLVM IRが標準出力に出力されます。  
        ```bash
        ./bfc hoge.bf -o hoge.ll
        ```
        オプションとして、LLVMの最適化のレベルを`-O0`～`-O3`で指定することが可能です。デフォルトは`-O3`で、レベルごとの最適化の内容はLLVMの各種コマンドと同じです。  
        ```bash
        ./bfc -O1 hoge.bf -o hoge.ll
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
   - Cコンパイラまで自動実行する方法
      1. PATHの通った場所に`llc`の名前でLLVM IRのコンパイラを、`cc`の名前でCコンパイラを用意します(コマンドエイリアス等では正しく動作しません)。  
      1. `-c`または`--compile-to-exe`フラグを付けてコンパイルします。  
         `-o`オプションの値は`cc`の`-o`オプションに直接渡されます(お使いの`cc`コマンドの仕様に注意してください)。`-o`オプションを省略すると、ソースファイルのパスから拡張子を消したパスを`cc`の`-o`オプションに渡します。  
         ```bash
         # -oオプションを省略した場合は`-o hoge`と同じ
         ./bfc --compile-to-exe hoge.bf -o hoge
         ```
1. 実行可能ファイルを実行します。  
   ```bash
   ./hoge
   ```

## License
ソースコードはMIT Licenseの下に配布されています。  
詳しくは[ライセンス表示](./LICENSE)や[https://www.tldrlegal.com/license/mit-license](https://www.tldrlegal.com/license/mit-license)などを参照>して下さい。  
なお、ライセンスは予告無く変更される可能性があり、変更前のライセンスは、そのライセンス表示がこのリポジトリのmainブランチの最新コミットに適用されていた期間に公開されたファイルにのみ適用されます。  
