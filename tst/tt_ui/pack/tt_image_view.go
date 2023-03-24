package pack

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func TTImageView(args *TTImageViewArgs) Widget {
	return ImageView{
		AssignTo: args.assignTo,
		Image:    args.image,
		MinSize:  args.customSize,
		MaxSize:  args.customSize,
		Mode:     args.mode,
	}
}

func NewTTImageViewArgs(assignTo **walk.ImageView) *TTImageViewArgs {
	return &TTImageViewArgs{assignTo: assignTo, mode: ImageViewModeShrink}
}

type TTImageViewArgs struct {
	assignTo   **walk.ImageView
	image      Property
	customSize Size
	mode       ImageViewMode
}

func (customT *TTImageViewArgs) SetAssignTo(assignTo **walk.ImageView) *TTImageViewArgs {
	customT.assignTo = assignTo
	return customT
}

func (customT *TTImageViewArgs) SetImage(image Property) *TTImageViewArgs {
	customT.image = image
	return customT
}

func (customT *TTImageViewArgs) SetCustomSize(customSize Size) *TTImageViewArgs {
	customT.customSize = customSize
	return customT
}

func (customT *TTImageViewArgs) SetMode(mode ImageViewMode) *TTImageViewArgs {
	customT.mode = mode
	return customT
}
