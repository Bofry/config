package resource

import (
	"fmt"
	"strings"

	"github.com/Bofry/structproto"
)

var _ structproto.TagResolver = ResourceTagResolver

func ResourceTagResolver(fieldname, token string) (*structproto.Tag, error) {
	if len(token) > 0 {
		parts := strings.SplitN(token, ";", 2)
		var desc string
		if len(parts) == 2 {
			parts, desc = strings.Split(parts[0], ","), parts[1]
		} else {
			parts = strings.Split(token, ",")
		}
		name, flags := parts[0], parts[1:]

		if len(flags) == 0 {
			for ii := 0; ii < len(name); ii++ {
				ch := name[ii]
				switch ch {
				case '*':
					flags = append(flags, structproto.RequiredFlag)
					continue
				}

				name = name[ii:]
				break
			}
		}

		if strings.HasSuffix(name, ".") {
			return nil, fmt.Errorf("invalid period symbol at end in resource name")
		}
		if strings.HasSuffix(name, " ") {
			return nil, fmt.Errorf("invalid space at end in resource name")
		}
		if strings.HasPrefix(name, " ") {
			return nil, fmt.Errorf("invalid space at start in resource name")
		}

		var tag *structproto.Tag
		if len(name) > 0 && name != "-" {
			tag = &structproto.Tag{
				Name:  name,
				Flags: flags,
				Desc:  desc,
			}
		}
		return tag, nil
	}
	return nil, nil
}
