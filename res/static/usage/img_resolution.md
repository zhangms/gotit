----------------------------------------------------
STEP1: 创建文件夹 workspace
STEP2: 在workspace下新建文本文档 resolution.txt，写入目标分辨率，可以换行写入多个分辨率
举个例子：文件内有三行

300*300
100*0
0*200

代表要将图片生成3种分辨率，分别为
  第一行：将图片宽高调整为300*300
  第二行：将图片宽度调整为100，高度按比例缩放
  第三行：将图片高度调整为200，宽度按比例缩放

STEP4: 在workspace下新建文件夹 resolution, 将需要调整分辨率的图片或文件夹放入其中
STEP5: 将本程序放入 workspace，双击运行，调整分辨率后的文件夹将放入 resolution_new 内

最终目录结构如下：

workspace
|-resolution.txt
|-gotit.exe
|-resolution
  |-1.jpg
  |-2.png
  |-dir1
    |-a.jpg
    |-b.jpeg
  |-dir2
    |-c.png
    |-d.PNG
----------------------------------------------------