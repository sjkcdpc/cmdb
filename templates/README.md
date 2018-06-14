curl -X POST http://127.0.0.1:8130/api/v1/cmdb/hosts -d "hostname=gate&conf=2core*4G*20M*50G+300G&wanip=182.254.145.181&lanip=10.0.0.11&os=Tencent Linux Release 2.2 (Final)&contact=madongsheng&manager=madongsheng&" | python -m json.tool

		Wanip:    wanip,
		Lanip:    lanip,
		Conf:     conf,
		Hostname: hostname,
		Os:       os,
		Contact:  contact,
		Manager:  manager,
		Tags:     tags,
		Remark:   remark,