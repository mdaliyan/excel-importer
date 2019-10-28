package inflxdb

func Write(tags map[string]string, fields map[string]interface{}) {

	if bp, err := NewBatchPoints(); err == nil {

		point, err := NewStatisticPoint(tags, fields)

		bp.AddPoint(point)

		if err := Client.Write(bp); err != nil {
			fn.LogError(err)
		}
	}

}
