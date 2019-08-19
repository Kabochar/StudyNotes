# 各种加密算法在Go语言中的使用

来源：<https://blog.51cto.com/634435/2117647>

# 使用SHA256、MD5、RIPEMD160

```
import (
    "fmt"
    "crypto/sha256"
    "os"
    "io"
    "crypto/md5"
    "golang.org/x/crypto/ripemd160"
)

func main()  {
    str := "hello world"
    sum := sha256.Sum256([]byte(str))
    fmt.Printf("SHA256：%x\n", sum)

    fileSha156()

    result := md5.Sum([]byte(str))
    fmt.Printf("MD5：%x\n", result)

    hasher := ripemd160.New()
    // 将加密内容的字节数组拷贝到ripemd160
    hasher.Write([]byte(str))
    fmt.Printf("RIPEMD160：%x", hasher.Sum(nil))
}

/**
 * 使用SHA256加密文件内容
 */
func fileSha156() {
    file, err := os.OpenFile("e:/test.txt", os.O_RDONLY, 0777)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    h := sha256.New()
    // 将文件内容拷贝到sha256中
    io.Copy(h, file)
    fmt.Printf("%x\n", h.Sum(nil))
}
```

# 使用DES

```
import (
    "bytes"
    "crypto/cipher" //cipher密码
    "crypto/des"
    "encoding/base64" //将对象转换成字符串
    "fmt"
)

/**
 * DES加密方法
 */
func MyDesEncrypt(orig, key string) string{

    // 将加密内容和秘钥转成字节数组
    origData := []byte(orig)
    k := []byte(key)

    // 秘钥分组
    block, _ := des.NewCipher(k)

    //将明文按秘钥的长度做补全操作
    origData = PKCS5Padding(origData, block.BlockSize())

    //设置加密方式－CBC
    blockMode := cipher.NewCBCDecrypter(block, k)

    //创建明文长度的字节数组
    crypted := make([]byte, len(origData))

    //加密明文
    blockMode.CryptBlocks(crypted, origData)

    //将字节数组转换成字符串，base64编码
    return base64.StdEncoding.EncodeToString(crypted)

}

/**
 * DES解密方法
 */
func MyDESDecrypt(data string, key string) string {

    k := []byte(key)

    //将加密字符串用base64转换成字节数组
    crypted, _ := base64.StdEncoding.DecodeString(data)

    //将字节秘钥转换成block快
    block, _ := des.NewCipher(k)

    //设置解密方式－CBC
    blockMode := cipher.NewCBCEncrypter(block, k)

    //创建密文大小的数组变量
    origData := make([]byte, len(crypted))

    //解密密文到数组origData中
    blockMode.CryptBlocks(origData, crypted)

    //去掉加密时补全的部分
    origData = PKCS5UnPadding(origData)

    return string(origData)
}

/**
 * 实现明文的补全
 * 如果ciphertext的长度为blockSize的整数倍，则不需要补全
 * 否则差几个则被几个，例：差5个则补5个5
 */
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

/**
 * 实现去补码，PKCS5Padding的反函数
 */
func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func main() {

    orig := "Hello World!"
    fmt.Println("原文：", orig)

    //声明秘钥,利用此秘钥实现明文的加密和密文的解密，长度必须为8
    key := "12345678"

    //加密
    encyptCode := MyDesEncrypt(orig, key)
    fmt.Println("密文：", encyptCode)

    //解密
    decyptCode := MyDESDecrypt(encyptCode, key)
    fmt.Println("解密结果：", decyptCode)
}
```

# 使用3DES

```
import (
    "bytes"
    "crypto/cipher"
    "crypto/des"
    "encoding/base64"
    "fmt"
)

func main() {
    orig := "hello world"
    // 3DES的秘钥长度必须为24位
    key := "123456781234567812345678"
    fmt.Println("原文：", orig)

    encryptCode := TripleDesEncrypt(orig, key)
    fmt.Println("密文：", encryptCode)

    decryptCode := TipleDesDecrypt(encryptCode, key)
    fmt.Println("解密结果：", decryptCode)

}

/**
 * 加密
 */
func TripleDesEncrypt(orig, key string) string {
    // 转成字节数组
    origData := []byte(orig)
    k := []byte(key)

    // 3DES的秘钥长度必须为24位
    block, _ := des.NewTripleDESCipher(k)
    // 补全码
    origData = PKCS5Padding(origData, block.BlockSize())
    // 设置加密方式
    blockMode := cipher.NewCBCEncrypter(block, k[:8])
    // 创建密文数组
    crypted := make([]byte, len(origData))
    // 加密
    blockMode.CryptBlocks(crypted, origData)

    return base64.StdEncoding.EncodeToString(crypted)
}

/**
 * 解密
 */
func TipleDesDecrypt(crypted string, key string) string {
    // 用base64转成字节数组
    cryptedByte, _ := base64.StdEncoding.DecodeString(crypted)
    // key转成字节数组
    k := []byte(key)

    block, _ := des.NewTripleDESCipher(k)
    blockMode := cipher.NewCBCDecrypter(block, k[:8])
    origData := make([]byte, len(cryptedByte))
    blockMode.CryptBlocks(origData, cryptedByte)
    origData = PKCS5UnPadding(origData)

    return string(origData)
}

func PKCS5Padding(orig []byte, size int) []byte {
    length := len(orig)
    padding := size - length%size
    paddintText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(orig, paddintText...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
```

# 使用AES

```
import (
    "bytes"
    "crypto/aes"
    "fmt"
    "crypto/cipher"
    "encoding/base64"
)

func main() {
    orig := "hello world"
    key := "123456781234567812345678"
    fmt.Println("原文：", orig)

    encryptCode := AesEncrypt(orig, key)
    fmt.Println("密文：" , encryptCode)

    decryptCode := AesDecrypt(encryptCode, key)
    fmt.Println("解密结果：", decryptCode)
}

func AesEncrypt(orig string, key string) string {
    // 转成字节数组
    origData := []byte(orig)
    k := []byte(key)

    // 分组秘钥
    block, _ := aes.NewCipher(k)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 补全码
    origData = PKCS7Padding(origData, blockSize)
    // 加密模式
    blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
    // 创建数组
    cryted := make([]byte, len(origData))
    // 加密
    blockMode.CryptBlocks(cryted, origData)

    return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, key string) string {
    // 转成字节数组
    crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
    k := []byte(key)

    // 分组秘钥
    block, _ := aes.NewCipher(k)
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 加密模式
    blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
    // 创建数组
    orig := make([]byte, len(crytedByte))
    // 解密
    blockMode.CryptBlocks(orig, crytedByte)
    // 去补全码
    orig = PKCS7UnPadding(orig)
    return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
    padding := blocksize - len(ciphertext)%blocksize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
```