---
title: Python实现身份证号码校验：结构与算法解析
date: 2019-10-01
tags: [Python, 算法, 工具]
categories: [Python]
excerpt: 解析中国二代身份证号码的结构组成，给出基于ISO 7064:1983.MOD 11-2校验码算法的Python实现，同时附带JavaScript版本供前端使用。
---

## 二代身份证号码结构

中国大陆二代身份证号码共 18 位，各段含义如下：

- **1-6 位**：行政区划代码
  - 1、2 位：所在省（直辖市、自治区）代码
  - 3、4 位：所在地级市（自治州）代码
  - 5、6 位：所在区（县、自治县、县级市）代码
- **7-14 位**：出生年月日（如 `19900101`）
- **15-16 位**：所在地派出所代码
- **第 17 位**：性别标识，奇数（1、3、5、7、9）为男性，偶数（2、4、6、8、0）为女性
- **第 18 位**：校验位，取值为 `0, 1, 2, 3, 4, 5, 6, 7, 8, 9, X`，由前 17 位按固定公式计算得出

第十八位的校验码依据 **ISO 7064:1983.MOD 11-2** 算法计算得出。

## JavaScript 实现

函数参数必须是**字符串**，因为二代身份证号码是 18 位，而在 JavaScript 中 18 位数值会超出计算精度范围，导致最后两位与计算的值不一致。

```javascript
function checkIDCard(idcode) {
    // 加权因子
    var weight_factor = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2];

    // 校验码对应表
    var check_code = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'];
    var code = idcode + "";
    var last = idcode[17]; // 最后一位

    var seventeen = code.substring(0, 17);

    // ISO 7064:1983.MOD 11-2
    // 判断最后一位校验码是否正确
    var arr = seventeen.split("");
    var len = arr.length;
    var num = 0;
    for (var i = 0; i < len; i++) {
        num = num + arr[i] * weight_factor[i];
    }

    // 获取余数
    var resisue = num % 11;
    var last_no = check_code[resisue];

    // 格式正则
    // 第一位不可能是0
    // 第二位到第六位可以是0-9
    // 第七位到第十位是年份，七八位为19或者20
    // 十一位和十二位是月份，01-12
    // 十三位和十四位是日期，01-31
    // 十五、十六、十七位是数字0-9
    // 十八位可能是数字0-9，也可能是X
    var idcard_patter = /^[1-9][0-9]{5}([1][9][0-9]{2}|[2][0][0|1][0-9])([0][1-9]|[1][0|1|2])([0][1-9]|[1|2][0-9]|[3][0|1])[0-9]{3}([0-9]|[X])$/;

    // 判断格式是否正确
    var format = idcard_patter.test(idcode);

    // 校验码和格式同时正确才是合法的身份证号码
    return last === last_no && format ? true : false;
}
```

## Python 实现

```python
# -*- coding: utf-8 -*-
def get_check_number(idcard):
    """
    根据前17位计算校验码
    :param idcard: 18位身份证号码（字符串）
    :return: 校验码字符
    """
    number = idcard[:17]
    xi_list = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]  # 每个位上乘的系数列表
    check_number = ["1", "0", "x", "9", "8", "7", "6", "5", "4", "3", "2"]  # 校验码对应表
    sum_of_number = 0
    for index in range(len(number)):
        sum_of_number += int(number[index]) * xi_list[index]
    # 取余数
    mod_num = sum_of_number % 11
    return check_number[mod_num]

print(get_check_number("4XXXXXXXXXXXXXXXXX"))
```

> 参考文章：[小雨编程](https://cloud.tencent.com/developer/article/1690115)
