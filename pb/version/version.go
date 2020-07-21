package version

import (
	"fmt"
	"keywea.com/cloud/pblib/pbconfig"
)

const buildPrefix = " build-"

// 主版本号 . 子版本号 [. 修正版本号 build- [编译版本号 ]]
// eg: 1.2.0 build-1234
type PBVersion struct {
	major int
	minor int
	revision int // 修正版本号
	buildver string // 编译版本号
	preRelease string // alpha(内部) / beta(公测)
}

func NewPBVersion(configor pbconfig.Configor) (*PBVersion, error) {
	if configor == nil {
		return &PBVersion{}, nil
	}
	major, _ := configor.GetInt("major", 0)
	minor, _ := configor.GetInt("minor", 0)
	revision, _ := configor.GetInt("revision", 0)

	return &PBVersion{
		major: major,
		minor: minor,
		revision: revision,
		buildver: configor.GetString("buildver", ""),
		preRelease: configor.GetString("preRelease", ""),
	}, nil
}

func (v *PBVersion) Alpha() {
	v.preRelease = "alpha"
}

func (v *PBVersion) Beta() {
	v.preRelease = "beta"
}

func (v *PBVersion) String() string {
	prerelease := v.preRelease;
	if len(prerelease) > 0 {
		prerelease = " " + prerelease
	}
	return fmt.Sprintf("%d.%d.%d%s%s", v.major, v.minor, v.revision, buildPrefix+v.buildver, prerelease)
}

func (v *PBVersion) After(b *PBVersion) bool {
	if v.major > b.major {
		return true
	}
	if v.minor > b.minor {
		return true
	}
	if v.revision > b.revision {
		return true
	}
	return false
}

func (v *PBVersion) Equal(b *PBVersion) bool {
	if v.major != b.major {
		return false
	}
	if v.minor != b.minor {
		return false
	}
	if v.revision != b.revision {
		return false
	}
	return true
}