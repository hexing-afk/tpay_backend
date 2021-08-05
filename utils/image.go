package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"path"
	"strings"
)

type UploadImageData struct {
	FileName     string //名称
	ContentType  string //图片类型
	Ext          string //文件后缀
	ImageContent []byte
}

const ImagPath = "image"

//图片base64字符串
//文件夹名称（存到哪个文件夹下）
func DealWithImage(imageStr, preName string) (*UploadImageData, error) {
	point := GetDailyId()

	if strings.HasPrefix(imageStr, "data:image/jpeg;base64,") {
		s := strings.Split(imageStr, "data:image/jpeg;base64,")
		imageStr = s[1]
	} else if strings.HasPrefix(imageStr, "data:image/png;base64,") {
		s := strings.Split(imageStr, "data:image/png;base64,")
		imageStr = s[1]
	}

	// 对文件名进行hash
	var baseName string
	img, err := base64.StdEncoding.DecodeString(imageStr)
	if err != nil {
		return nil, err
	}

	// 获取文件后缀
	ext := GetFileTypeFromMagic(img)
	switch ext {
	case "png":
		baseName = Md5(point) + ".png"
	case "jpeg":
		baseName = Md5(point) + ".jpeg"
	default:
		return nil, errors.New(fmt.Sprintf("上传的图片格式不被支持:%s", ext))
	}

	// 重新生成文件，防止文件内容中有攻击代码
	newImg, _, err := RecreateImage(img)
	if len(newImg) == 0 {
		return nil, errors.New(fmt.Sprintf("重新生成文件失败, error: %v", err))
	}

	fileName := path.Join(preName, ImagPath, baseName)
	contentType := http.DetectContentType(newImg)

	data := &UploadImageData{
		ContentType:  contentType,
		FileName:     fileName,
		Ext:          ext,
		ImageContent: newImg,
	}

	return data, nil
}

/*
 * 解析base64编码格式图片
 * 返回: 图片后缀名,图片二进制内容,错误
 * 参数格式示例:
 * jpeg: data:image/jpeg;base64,/9j/4AAQSkZJRg...
 * png: data:image/png;base64,iVBORw0KGgo...
 */
func ParseBase64Img(imgBase64Str string) (ext string, data []byte, err error) {
	arr := strings.Split(imgBase64Str, ";base64,")
	if len(arr) != 2 || arr[0] == "" || arr[1] == "" {
		return "", nil, errors.New("The content is incorrect")
	}

	ext = strings.TrimPrefix(arr[0], "data:image/")
	if ext == "" {
		return "", nil, errors.New("Failed to get suffix")
	}

	data, err = base64.StdEncoding.DecodeString(arr[1])
	if err != nil {
		return "", nil, err
	}

	return ext, data, nil
}

func RecreateImage(img []byte) (newImg []byte, ext string, e error) {
	imgS, ext, err := image.Decode(bufio.NewReader(bytes.NewReader(img)))
	if imgS == nil {
		return nil, ext, err
	}
	//fmt.Println(imgS)

	x := imgS.Bounds().Size().X
	y := imgS.Bounds().Size().Y
	fmt.Printf("x=[%v],y=[%v],ext=[%v]\n", x, y, ext)
	m := image.NewRGBA(imgS.Bounds())
	draw.Draw(m, imgS.Bounds(), imgS, image.ZP, draw.Src)
	pngDst := bytes.NewBufferString("")
	switch ext {
	case "png":
		err = png.Encode(pngDst, m)
	case "jpeg":
		err = jpeg.Encode(pngDst, m, &jpeg.Options{90})
	default:
		return nil, ext, errors.New("only support png && jpeg")
	}
	if err != nil {
		return nil, ext, err
	}
	newImg = pngDst.Bytes()
	return newImg, ext, nil
}

func GetFileTypeFromMagic(img []byte) (fileType string) {
	magics := map[string]string{
		//WMV
		"png":       `89504E47`,
		"jpeg":      `FFD8FF`,
		"gif":       `47494638`,
		"tiff":      `49492A00`,
		"bmp":       "424D",
		"pbm_a":     "5031",
		"pgm_a":     "5032",
		"ppm_a":     "5033",
		"pbm_b":     "5034",
		"pgm_b":     "5035",
		"ppm_pnm_b": "5036",
		"dwg":       "41433130",
		"psd":       "38425053",
		"rtf":       "7B5C727466",
		"xml":       "3C3F786D6C",
		"html":      "68746D6C3E",
		"eml":       "44656C69766572792D646174653A",
		"dbx":       "CFAD12FEC5FD746F",
		"pst":       "2142444E",
		"xls/doc":   "D0CF11E0",
		"mdb":       "5374616E64617264204A",
		"wpd":       "FF575043",
		"eps/ps":    "252150532D41646F6265",
		"pdf":       "255044462D312E",
		"qdf":       "AC9EBD8F",
		"pwl":       "E3828596",
		"zip":       "504B0304",
		"rar":       "52617221",
		"7z":        "377ABCAF271C",
		"bz2":       "425A",
		"gz":        "1F8B",
		"xz":        "FD377A585A00",
		"llvm":      "4243",
		"wav":       "57415645",
		"avi":       "41564920",
		"ram":       "2E7261FD",
		"rm":        "2E524D46",
		"mpg":       "000001BA",
		"mpg_2":     "000001B3",
		"mov":       "6D6F6F76",
		"asf":       "3026B2758E66CF11",
		"mid":       "4D546864",
		// exec
		"a.out":   `0107`,
		"java":    `CAFEBABE`,
		"pe/coef": `4d5a`,
		"elf":     `7F454C46`,
		"script":  `2321`,
		"mp4":     `00000020667`,
		"3gp":     `0000001c66747970`,
	}
	fileType = ""
	for k, v := range magics {
		isMagic := chkMagic(img, v)
		if isMagic {
			fmt.Printf("is[%v]\n", k)
			fileType = k
			break
		}
	}
	if fileType == "" {
		fmt.Printf("maybe [plain text]or[*.tar]\n")
	}
	return fileType
}

func chkMagic(img []byte, magic string) bool {
	b, _ := hex.DecodeString(magic)
	isMagic := bytes.HasPrefix(img, b)
	return isMagic
}
