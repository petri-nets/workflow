package flow

// {
//     "name": "xxxx工作流",
//     "if": [ // 顺序判断，
//         {
//            "name": "作业1",
//            "jobID": "123",
//            "condition": "TDJOB_JOB_ID > 1"
//         },
//         {
//             "name": "作业2",
//             "jobID": "234",
//             "condition": "TDJOB_WORKITEM > 2"
//         }
//     ]
// }

// {
//     "name": "xxxx工作流",
//     "sequel": [ // 顺序执行，前面失败，后面的继续执行
//         {
//            "name": "作业1",
//            "jobID": "123",
//         },
//         {
//             "name": "作业2",
//             "jobID": "234",
//         },
//     ]
// }

// {
//     "name": "xxxx工作流",
//     "serial": [ // 串行执行，前面执行失败，后面将终止
//         {
//            "name": "作业1",
//            "jobID": "123",
//         },
//         {
//             "name": "作业2",
//             "jobID": "234",
//         },
//     ]
// }

// {
//     "name": "xxxx工作流",
//     "parallel": [ // 并行执行，前面执行失败，后面将终止
//         {
//            "name": "作业1",
//            "jobID": "123",
//         },
//         {
//             "name": "作业2",
//             "jobID": "234",
//         },
//     ]
// }

// sequel(seq), serial(ser), parallel(par), then, if, elseif

// {
//     "name": "xxxx工作流",
//     "sequel": [ // 并行执行，前面执行失败，后面将终止
//         {
//            "name": "作业1",
//            "jobID": "123",
//         },
//         {
//             "name": "作业2",
//             "jobID": "234",
//         },
//     ],
// 	"then": {	// then 表示等待前面所有满足条件的执行完成后，才开始执行
// 		"name": "作业3",
// 		"jobID": "456",
// 	}
// }

// 工作流基本信息
// {
// 	"id": 0, // 工作流ID，如果为0，则新增，否则修改
// 	"name": "工作流1", // 工作流名 require
// 	"desc": "", // 描述
// 	"start_job_id": "123",
// 	"start_date": "0000-00-00 00:00:00", // 开始启用使用，默认 0000-00-00 00:00:00 表示立即开始，否则按指定时间开始
// 	"end_date": "", // 可以不填写，默认无结束时间，否则根据指定结束时间，决定是否生效
// 	"operator": "", // 操作人
// }

// 工作流任务基本信息
// {
// 	"name": "作业1", // 工作流名 require
// 	"desc": "", // 描述
// 	"trigger": "AUTO", // 任务触发方式，AUTO（自动触发），USER（用户触发），TIMER（定时器触发），MSG（消息）
// 	"job": "123", // 作业ID
// 	"role_id": "12", // 角色ID，用于限制允许操作的用户组，针对trigger为USER类型的才有效
// }

// 任务可以层层嵌套

// {
// 	"name": "工作流1",
// 	"desc": "",
// 	"start_job_id": "123",
// 	"start_date": "0000-00-00 00:00:00",
// 	"end_date": "",
// 	"operator": "edenzou",
//  "router": "seq",
// 	"tasks": [
// 		{
// 			"name": "作业1",
// 			"desc": "",
// 			"trigger": "AUTO",
// 			"job": "1",
// 			"role_id": "11",
// 			"par": [
// 				{
// 					"name": "作业11",
// 					"desc": "",
// 					"trigger": "AUTO",
// 					"job": "111",
// 					"role_id": "12",
// 				},
// 				{
// 					"name": "作业12",
// 					"desc": "",
// 					"trigger": "AUTO",
// 					"job": "112",
// 					"role_id": "12",
// 				},
// 			],
// 			"then": {
// 				"name": "作业",
// 				"desc": "",
// 				"trigger": "AUTO",
// 				"job": "12",
// 				"role_id": "12",
// 			},
// 		},
// 		{
// 			"name": "作业2",
// 			"desc": "",
// 			"trigger": "AUTO",
// 			"job": "2",
// 			"role_id": "12",
// 		}
// 	],
// 	"then": {
// 		"name": "作业3",
// 		"desc": "",
// 		"trigger": "AUTO",
// 		"job": "3",
// 		"role_id": "12",
// 	}
// }
