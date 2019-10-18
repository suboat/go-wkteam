package wkteam

// 消息类型转字符串
// // 消息类型：1文字、2图片、3表情、 4语音、5视频、6文件、10系统消息
func priContentTypeToStr(in int) (out string) {
	switch in {
	case 1:
		out = MsgCategoryTxt
	case 2:
		out = MsgCategoryImg
	case 3:
		out = MsgCategoryGif
	case 4:
		out = MsgCategoryWav
	case 5:
		out = MsgCategoryMp4
	case 6:
		out = MsgCategoryDoc
	case 10:
		out = MsgCategorySys
	}
	return
}
