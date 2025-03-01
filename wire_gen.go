// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"chatroom/config"
	"chatroom/controller"
	"chatroom/dao"
	"chatroom/dao/cache"
	"chatroom/pkg/business"
	"chatroom/pkg/client"
	"chatroom/pkg/core"
	socket2 "chatroom/pkg/core/socket"
	"chatroom/pkg/kafka"
	"chatroom/service"
	"chatroom/service/message"
	"chatroom/socket"
	"chatroom/socket/consume"
	chat2 "chatroom/socket/consume/chat"
	example2 "chatroom/socket/consume/example"
	"chatroom/socket/handler"
	"chatroom/socket/handler/event"
	"chatroom/socket/handler/event/chat"
	"chatroom/socket/handler/event/example"
	"chatroom/socket/process"
	"chatroom/socket/router"
	"github.com/google/wire"
)

// Injectors from wire.go:

func NewHttpInjector(conf *config.Config) *core.AppProvider {
	engine := core.NewGinServer()
	redisClient := client.NewRedisClient(conf)
	jwtTokenStorage := cache.NewTokenSessionStorage(redisClient)
	db := client.NewMySQLClient(conf)
	users := dao.NewUsers(db, redisClient)
	userService := &service.UserService{
		UsersRepo: users,
	}
	organize := dao.NewOrganize(db)
	user := &controller.User{
		Redis:        redisClient,
		Session:      jwtTokenStorage,
		Config:       conf,
		UserService:  userService,
		UsersRepo:    users,
		OrganizeRepo: organize,
	}
	admin := dao.NewAdmin(db)
	daoJwtTokenStorage := dao.NewTokenSessionStorage(redisClient)
	captchaStorage := dao.NewCaptchaStorage(redisClient)
	captcha := dao.NewBase64Captcha(captchaStorage)
	auth := &controller.Auth{
		Config:          conf,
		AdminRepo:       admin,
		UserRepo:        users,
		JwtTokenStorage: daoJwtTokenStorage,
		ICaptcha:        captcha,
		UserService:     userService,
	}
	redisLock := cache.NewRedisLock(redisClient)
	messageStorage := cache.NewMessageStorage(redisClient)
	serverStorage := cache.NewSidStorage(redisClient)
	clientStorage := cache.NewClientStorage(redisClient, conf, serverStorage)
	unreadStorage := cache.NewUnreadStorage(redisClient)
	contactRemark := cache.NewContactRemark(redisClient)
	relation := cache.NewRelation(redisClient)
	contact := dao.NewContact(db, contactRemark, relation)
	group := dao.NewGroup(db)
	source := dao.NewSource(db, redisClient)
	groupMember := dao.NewGroupMember(db, relation)
	talkService := &service.TalkService{
		Source:          source,
		GroupMemberRepo: groupMember,
	}
	vote := cache.NewVote(redisClient)
	groupVote := dao.NewGroupVote(db, vote)
	talkUserMessage := dao.NewTalkRecordFriend(db)
	talkGroupMessage := dao.NewTalkRecordGroup(db)
	talkGroupMessageDel := dao.NewTalkRecordGroupDel(db)
	talkRecordService := &service.TalkRecordService{
		Source:                source,
		TalkVoteCache:         vote,
		TalkRecordsVoteRepo:   groupVote,
		GroupMemberRepo:       groupMember,
		TalkRecordFriendRepo:  talkUserMessage,
		TalkRecordGroupRepo:   talkGroupMessage,
		TalkRecordsDeleteRepo: talkGroupMessageDel,
	}
	talkSession := dao.NewTalkSession(db)
	talkSessionService := &service.TalkSessionService{
		Source:          source,
		TalkSessionRepo: talkSession,
	}
	sequence := cache.NewSequence(redisClient)
	daoSequence := dao.NewSequence(db, sequence)
	kafkaClient := kafka.NewKafkaClient(conf)
	pushMessage := &business.PushMessage{
		Redis: redisClient,
		Kafka: kafkaClient,
	}
	groupService := &service.GroupService{
		Source:          source,
		GroupRepo:       group,
		GroupMemberRepo: groupMember,
		Relation:        relation,
		Sequence:        daoSequence,
		PushMessage:     pushMessage,
	}
	authService := &service.AuthService{
		OrganizeRepo:    organize,
		ContactRepo:     contact,
		GroupRepo:       group,
		GroupMemberRepo: groupMember,
	}
	contactService := &service.ContactService{
		Source:      source,
		ContactRepo: contact,
	}
	clientConnectService := &service.ClientConnectService{
		Storage: clientStorage,
	}
	session := &controller.Session{
		RedisLock:            redisLock,
		Session:              jwtTokenStorage,
		Config:               conf,
		MessageStorage:       messageStorage,
		ClientStorage:        clientStorage,
		UnreadStorage:        unreadStorage,
		ContactRemark:        contactRemark,
		ContactRepo:          contact,
		UsersRepo:            users,
		GroupRepo:            group,
		TalkService:          talkService,
		TalkRecordsService:   talkRecordService,
		TalkSessionService:   talkSessionService,
		UserService:          userService,
		GroupService:         groupService,
		AuthService:          authService,
		ContactService:       contactService,
		ClientConnectService: clientConnectService,
	}
	contactGroup := dao.NewContactGroup(db)
	contactGroupService := &service.ContactGroupService{
		Source:           source,
		ContactGroupRepo: contactGroup,
	}
	contactApplyService := &service.ContactApplyService{
		Source:      source,
		PushMessage: pushMessage,
	}
	fileUpload := dao.NewFileUpload(db)
	iFilesystem := config.NewFilesystem(conf)
	robot := dao.NewRobot(db)
	messageService := &message.Service{
		Source:              source,
		GroupMemberRepo:     groupMember,
		SplitUploadRepo:     fileUpload,
		TalkRecordsVoteRepo: groupVote,
		UsersRepo:           users,
		Filesystem:          iFilesystem,
		UnreadStorage:       unreadStorage,
		MessageStorage:      messageStorage,
		ServerStorage:       serverStorage,
		ClientStorage:       clientStorage,
		Sequence:            daoSequence,
		RobotRepo:           robot,
		Kafka:               kafkaClient,
		PushMessage:         pushMessage,
	}
	controllerContact := &controller.Contact{
		Session:             jwtTokenStorage,
		Config:              conf,
		ClientStorage:       clientStorage,
		ContactRepo:         contact,
		UsersRepo:           users,
		OrganizeRepo:        organize,
		ContactGroupRepo:    contactGroup,
		TalkSessionRepo:     talkSession,
		ContactService:      contactService,
		UserService:         userService,
		TalkListService:     talkSessionService,
		ContactGroupService: contactGroupService,
		ContactApplyService: contactApplyService,
		MessageService:      messageService,
	}
	groupApplyStorage := cache.NewGroupApplyStorage(redisClient)
	groupNotice := dao.NewGroupNotice(db)
	groupMemberService := &service.GroupMemberService{
		Source:          source,
		GroupMemberRepo: groupMember,
	}
	controllerGroup := &controller.Group{
		Session:            jwtTokenStorage,
		Config:             conf,
		RedisLock:          redisLock,
		Repo:               source,
		UsersRepo:          users,
		GroupRepo:          group,
		GroupMemberRepo:    groupMember,
		GroupApplyStorage:  groupApplyStorage,
		GroupNoticeRepo:    groupNotice,
		TalkSessionRepo:    talkSession,
		GroupService:       groupService,
		GroupMemberService: groupMemberService,
		TalkSessionService: talkSessionService,
		UserService:        userService,
		ContactService:     contactService,
	}
	emoticon := dao.NewEmoticon(db)
	emoticonService := &service.EmoticonService{
		Source:       source,
		EmoticonRepo: emoticon,
		Filesystem:   iFilesystem,
	}
	controllerEmoticon := &controller.Emoticon{
		Session:         jwtTokenStorage,
		Config:          conf,
		RedisLock:       redisLock,
		EmoticonRepo:    emoticon,
		EmoticonService: emoticonService,
		Filesystem:      iFilesystem,
	}
	publish := &controller.Publish{
		Session:        jwtTokenStorage,
		Config:         conf,
		AuthService:    authService,
		MessageService: messageService,
		Kafka:          kafkaClient,
	}
	fileSplitUploadService := &service.FileSplitUploadService{
		Source:          source,
		SplitUploadRepo: fileUpload,
		Config:          conf,
		FileSystem:      iFilesystem,
	}
	upload := &controller.Upload{
		Config:             conf,
		Filesystem:         iFilesystem,
		SplitUploadService: fileSplitUploadService,
		Session:            jwtTokenStorage,
	}
	controllers := &controller.Controllers{
		User:     user,
		Auth:     auth,
		Session:  session,
		Contact:  controllerContact,
		Group:    controllerGroup,
		Emoticon: controllerEmoticon,
		Publish:  publish,
		Upload:   upload,
	}
	appProvider := &core.AppProvider{
		Config:      conf,
		Engine:      engine,
		Controllers: controllers,
	}
	return appProvider
}

func NewSocketInjector(conf *config.Config) *socket.AppProvider {
	redisClient := client.NewRedisClient(conf)
	serverStorage := cache.NewSidStorage(redisClient)
	clientStorage := cache.NewClientStorage(redisClient, conf, serverStorage)
	clientConnectService := &service.ClientConnectService{
		Storage: clientStorage,
	}
	db := client.NewMySQLClient(conf)
	relation := cache.NewRelation(redisClient)
	groupMember := dao.NewGroupMember(db, relation)
	source := dao.NewSource(db, redisClient)
	groupMemberService := &service.GroupMemberService{
		Source:          source,
		GroupMemberRepo: groupMember,
	}
	kafkaClient := kafka.NewKafkaClient(conf)
	pushMessage := &business.PushMessage{
		Redis: redisClient,
		Kafka: kafkaClient,
	}
	chatHandler := &chat.Handler{
		Redis:         redisClient,
		Source:        source,
		MemberService: groupMemberService,
		PushMessage:   pushMessage,
	}
	roomStorage := socket2.NewRoomStorage()
	chatEvent := &event.ChatEvent{
		Redis:           redisClient,
		GroupMemberRepo: groupMember,
		MemberService:   groupMemberService,
		Handler:         chatHandler,
		RoomStorage:     roomStorage,
		PushMessage:     pushMessage,
	}
	chatChannel := &handler.ChatChannel{
		Storage: clientConnectService,
		Event:   chatEvent,
	}
	exampleHandler := example.NewHandler()
	exampleEvent := &event.ExampleEvent{
		Handler: exampleHandler,
	}
	exampleChannel := &handler.ExampleChannel{
		Storage: clientStorage,
		Event:   exampleEvent,
	}
	handlerHandler := &handler.Handler{
		Chat:        chatChannel,
		Example:     exampleChannel,
		Config:      conf,
		RoomStorage: roomStorage,
	}
	jwtTokenStorage := cache.NewTokenSessionStorage(redisClient)
	engine := router.NewRouter(conf, handlerHandler, jwtTokenStorage)
	healthSubscribe := process.NewHealthSubscribe(serverStorage)
	organize := dao.NewOrganize(db)
	users := dao.NewUsers(db, redisClient)
	vote := cache.NewVote(redisClient)
	groupVote := dao.NewGroupVote(db, vote)
	talkUserMessage := dao.NewTalkRecordFriend(db)
	talkGroupMessage := dao.NewTalkRecordGroup(db)
	talkGroupMessageDel := dao.NewTalkRecordGroupDel(db)
	talkRecordService := &service.TalkRecordService{
		Source:                source,
		TalkVoteCache:         vote,
		TalkRecordsVoteRepo:   groupVote,
		GroupMemberRepo:       groupMember,
		TalkRecordFriendRepo:  talkUserMessage,
		TalkRecordGroupRepo:   talkGroupMessage,
		TalkRecordsDeleteRepo: talkGroupMessageDel,
	}
	contactRemark := cache.NewContactRemark(redisClient)
	contact := dao.NewContact(db, contactRemark, relation)
	contactService := &service.ContactService{
		Source:      source,
		ContactRepo: contact,
	}
	handler2 := &chat2.Handler{
		Config:               conf,
		OrganizeRepo:         organize,
		UserRepo:             users,
		Source:               source,
		TalkRecordsService:   talkRecordService,
		ContactService:       contactService,
		ClientConnectService: clientConnectService,
		RoomStorage:          roomStorage,
	}
	chatSubscribe := consume.NewChatSubscribe(handler2)
	handler3 := example2.NewHandler()
	exampleSubscribe := consume.NewExampleSubscribe(handler3)
	messageSubscribe := process.NewMessageSubscribe(redisClient, chatSubscribe, exampleSubscribe, kafkaClient)
	subServers := &process.SubServers{
		HealthSubscribe:  healthSubscribe,
		MessageSubscribe: messageSubscribe,
	}
	server := process.NewServer(subServers)
	emailClient := client.NewEmailClient(conf)
	providers := &client.Providers{
		EmailClient: emailClient,
	}
	appProvider := &socket.AppProvider{
		Config:    conf,
		Engine:    engine,
		Coroutine: server,
		Handler:   handlerHandler,
		Providers: providers,
	}
	return appProvider
}

// wire.go:

var ProviderSet = wire.NewSet(client.NewMySQLClient, client.NewEmailClient, client.NewRedisClient, config.NewFilesystem, kafka.NewKafkaClient, wire.Struct(new(client.Providers), "*"))
