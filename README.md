# MunouLoki
最初に数十分で作られた人工無能です。

## 使い方 / Usage
### JA
1. Releasesから、使用している環境に合った「MunouLoki-NAME-amd64.zip」をクリックしてダウンロード
2. 解凍しフォルダーの奥へ歩みを進める
3. main.exeやmainファイルを発見するのでそれを実行 (linuxの場合は`chmod +x main`などで権限を与える必要がある)
4. うんたらかんたら
### EN
1. From the “Releases” section, click on the “MunouLoki-NAME-amd64.zip” file that matches your system to download it.
2. Unzip the file and navigate to the folder inside.
3. You'll find a file named “main.exe” or “main”; run it (on Linux, you'll need to grant execution permissions using a command like `chmod +x main`).
4. Blah, blah, blah...
## 仕組み / System
### JA
まず、記憶が空の状態でユーザーからメッセージを受け取るとオウム返しします。
そのオウム返しに対するユーザーの応答を観察し、記憶に保存します。
次回以降、ユーザーから入力された文字列と、自身が返した文字列との類似度を比較し、最も高い組み合わせのユーザー側の返答を用います。
(確率が同じ場合はランダムに選択される)
(十六分の一であえて無作為な返答をする)
### EN
First, when a message is received from the user while the memory is empty, the system parrots the message back.
It observes the user’s response to that parroted message and stores it in memory.
From then on, it compares the similarity between the string entered by the user and the strings it has previously returned, and uses the one with the highest similarity as its response.
(If the probabilities are the same, a response is selected at random.)
(One in ten times, it deliberately gives a random response.)
## その他 / Etc...
### JA
- main.exe実行時に警告が出る場合があります。基本的には大丈夫なんですが、実行が恐ろしいと思う人はリポジトリからmain.goをダウンロードして、GoLangコンパイラーで変換してから実行しても良い。
- memoryFileはただのJSONファイルです。
### EN
- A warning may appear when you run main.exe. It’s generally fine, but if you’re hesitant to run it, you can download main.go from the repository, compile it with the GoLang compiler, and then run it.
- memoryFile is just a JSON file.
