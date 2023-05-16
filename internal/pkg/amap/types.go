package amap

type DirectionResp struct {
	Status   string `json:"status"` // 0：请求失败；1：请求成功
	Info     string `json:"info"`   // status为0时，info返回错误原因，否则返回“OK”。
	Infocode string `json:"infocode"`
	Count    string `json:"count"` // 驾车路径规划方案数目
	Route    struct {
		Origin      string     `json:"origin"`      // 起点坐标
		Destination string     `json:"destination"` // 终点坐标
		Paths       []struct { // 驾车换乘方案(直接取第一个就行了)
			Distance     string     `json:"distance"` // 行驶距离(米)
			Duration     string     `json:"duration"` // 预计行驶时间(秒)
			Strategy     string     `json:"strategy"` // 路径策略
			Tolls        string     `json:"tolls"`    // 道路收费总金额(元)
			TollDistance string     `json:"toll_distance"`
			Steps        []struct { // 路径规划方案中的每一段路线
				Instruction     string `json:"instruction"`      // 行驶指示
				Orientation     string `json:"orientation"`      // 方向
				Distance        string `json:"distance"`         // 此路段距离
				Tolls           string `json:"tolls"`            // 此段收费
				TollDistance    string `json:"toll_distance"`    // 收费路段距离
				TollRoad        any    `json:"toll_road"`        // 主要收费道路
				Duration        string `json:"duration"`         // 此路段预计行驶时间
				Polyline        string `json:"polyline"`         // 此路段坐标点串
				Action          any    `json:"action"`           // 导航主要动作
				AssistantAction any    `json:"assistant_action"` // 导航辅助动作
				Road            string `json:"road,omitempty"`   // 道路名称
			} `json:"steps"`
			Restriction   string `json:"restriction"`    // 限行结果
			TrafficLights string `json:"traffic_lights"` // 红绿灯个数
		} `json:"paths"`
	} `json:"route"`
}

type DistrictResp struct {
	Status     string `json:"status"`
	Info       string `json:"info"`
	Infocode   string `json:"infocode"`
	Count      string `json:"count"`
	Suggestion struct {
		Keywords []string `json:"keywords"` // 建议关键字
		Cities   []string `json:"cities"`   // 建议城市
	} `json:"suggestion"`
	Districts []struct {
		Citycode  any    `json:"citycode"`
		Adcode    string `json:"adcode"`
		Name      string `json:"name"`
		Center    string `json:"center"`
		Level     string `json:"level"`
		Districts []struct {
			Citycode  any    `json:"citycode"`
			Adcode    string `json:"adcode"`
			Name      string `json:"name"`
			Center    string `json:"center"`
			Level     string `json:"level"`
			Districts []any  `json:"districts"` // 下级行政区列表，包含district元素
		} `json:"districts"`
	} `json:"districts"`
}
