----------------------------------------------------
STEP1: 创建文件夹 workspace
STEP2: 打开网址：https://tinify.cn/developers 获取API key
STEP3: 在workspace下新建文本文档 apikey.txt，将 STEP2 获取到的API key 按行写入 apikey.txt
STEP4: 在workspace下新建文件夹 compress, 将需要压缩的图片或文件夹放入其中
STEP5: 将本程序放入 workspace，双击运行，压缩后的文件夹将放入 compress_new 内
NOTE:  一个 tinypng 账号默认只能一个月免费压缩500张图片，超过500张可以付费或者换个邮箱获取apikey

最终目录结构如下：

workspace
|-apikey.txt
|-gotit.exe
|-compress
  |-1.jpg
  |-2.png
  |-dir1
    |-a.jpg
    |-b.jpeg
  |-dir2
    |-c.png
    |-d.PNG
----------------------------------------------------