# 随手记导出

无需充值会员即可将“随手记”APP中的记账记录通通导出为Excel表格。导出的Excel使用了[“钱迹”](https://qianjiapp.com/)的模板。

### 1. 下载可执行文件或自行编译

[https://github.com/rroy233/suishoujiExport/releases
](https://github.com/rroy233/suishoujiExport/releases)

### 2. 获取cookie

- 前往 https://www.sui.com/tally/new.do ，登录
- 按`F12`打开开发者工具
- 在`网络`选项卡中，找到`/tally/new.rmi`的POST请求，复制请求头的`cookie`

### 3. 配置

文件名为`config.json`。将需要导出账单的**起止日期**以及上面获得的**cookie**填入。

```json
{
  "cookie": "xxxxxxx",
  "beginTime": "2023.01.01",
  "endTime": "2023.12.31"
}
```

### 4. 运行

导出文件名为`out-起始日期-终止日期.xlsx`

参考钱迹的[导入教程](https://docs.qianjiapp.com/other/import_templete.html)进行导入。