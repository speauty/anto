package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func UIImageView(args *UIImageViewArgs) Widget {
	return ImageView{
		AssignTo: args.assignTo,
		Image:    args.image,
		MinSize:  args.customSize,
		MaxSize:  args.customSize,
		Mode:     args.mode,
	}
}

func NewUIImageViewArgs(assignTo **walk.ImageView) *UIImageViewArgs {
	return &UIImageViewArgs{assignTo: assignTo, mode: ImageViewModeShrink}
}

type UIImageViewArgs struct {
	assignTo   **walk.ImageView
	image      Property
	customSize Size
	mode       ImageViewMode
}

func (customT *UIImageViewArgs) SetAssignTo(assignTo **walk.ImageView) *UIImageViewArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *UIImageViewArgs) SetImage(image Property) *UIImageViewArgs {
	customT.image = image
	return customT
}

func (customT *UIImageViewArgs) SetCustomSize(customSize Size) *UIImageViewArgs {
	customT.customSize = customSize
	return customT
}

func (customT *UIImageViewArgs) SetMode(mode ImageViewMode) *UIImageViewArgs {
	customT.mode = mode
	return customT
}
