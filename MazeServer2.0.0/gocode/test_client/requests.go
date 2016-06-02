// test_requests
package main

import (
    //"strings"
    //"bufio"
    //"os"
    "runtime"
    "time"
    //"flag"
    //"fmt"

    //	proto "code.google.com/p/goprotobuf/proto"
    //beegoconf "github.com/astaxie/beego/config"
    xylog "guanghuan.com/xiaoyao/common/log"
    //"guanghuan.com/xiaoyao/superbman_server/money"
)

const YAML_CONF_FILE = "test.yaml"

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU())

    initConfig()
    //DB_Init()
    //xymoney.Init(DefDB)

    var err error

    //初始化一个goruntine用于定时输出测试进度
    done, id, sum, errSum := uint64(0), DefConfig.Id, DefConfig.Sum, uint64(0)
    go test_print_progress(&done, &errSum, sum)

    switch DefConfig.Cmd {
    //login
    case TEST_CMD_CODE_LOGIN:
        err = test_login()
    case TEST_CMD_CODE_PPROF_LOGIN:
        err = test_pprof_login(&done, &errSum, id, sum)
    //bind
    case TEST_CMD_CODE_BIND:
        err = test_bind()
    //case TEST_CMD_CODE_PPROF_BIND:
    //	err = test_pprof_bind(&done, &errSum, id, sum)
    //announcement
    case TEST_CMD_CODE_QUERY_ANNOUNCEMENT: // query user announcement
        err = test_announcement_querybytime()
    case TEST_CMD_CODE_PPROF_QUERY_ANNOUNCEMENT: // pprof query announcement
        err = test_pprof_announcement_querybytime(&done, &errSum, id, sum)
    //iap
    case TEST_CMD_CODE_IAP_VARIFY: // iapverify
        err = test_iap()
    case TEST_CMD_CODE_PPROF_IAP_VARIFY: // pprof iapverify
        //todo
    //goods
    case TEST_CMD_CODE_QUERY_GOODS: // good query
        err = test_query_goods()
    case TEST_CMD_CODE_PPROF_QUERY_GOODS: // pprof good query
        err = test_pprof_query_goods(&done, &errSum, id, sum)
    case TEST_CMD_CODE_BUY_GOODS: // good buy
        err = test_buy_good()
    case TEST_CMD_CODE_PPROF_BUY_GOODS: // pprof good buy
        err = test_pprof_buy_good(&done, &errSum, id, sum)
    //stamina
    //case TEST_CMD_CODE_QUERY_STAMINA: // query user stamina
    //	err = test_query_user_stamina()
    //frienddata
    case TEST_CMD_CODE_QUERY_FRIENDDATA: //friend data
        err = test_query_frienddata()
    case TEST_CMD_CODE_PPROF_QUERY_FRIENDDATA:
        err = test_pprof_query_frienddata(&done, &errSum, id, sum)
    //syslotto
    case TEST_CMD_CODE_SYSLOTTO: //syslotto
        err = test_sys_lotto()
    case TEST_CMD_CODE_PPROF_SYSLOTTO:
        err = test_pprof_sys_lotto(&done, &errSum, id, sum)
    //game
    case TEST_CMD_CODE_NEWGAME: //new game
        err = test_new_game()
    case TEST_CMD_CODE_PPROF_NEWGAME: // pprof new game
        err = test_pprof_new_game(&done, &errSum, id, sum)
    case TEST_CMD_CODE_PPROF_NEWGAMEANDRESULT:
        err = test_pprof_game(&done, &errSum, id, sum)
    case TEST_CMD_CODE_GAMERESULT: //game result
        err = test_game_result2()
    //mission
    case TEST_CMD_CODE_QUERY_USERMISSION: //query user mission
        err = test_query_user_mission()
    case TEST_CMD_CODE_PPROF_QUERY_USERMISSION: //pprof query user mission
        err = test_pprof_query_user_mission(&done, &errSum, id, sum)
    case TEST_CMD_CODE_CONFIRM_USERMISSION: //confirm user mission
        err = test_confirm_user_mission()
    case TEST_CMD_CODE_PPROF_CONFIRM_USERMISSION: //pprof confirm user mission
        err = test_pprof_confirm_user_mission(&done, &errSum, id, sum)
    //activity
    case TEST_CMD_CODE_QUERY_SIGNIN_ACTIVITY: //query user signin activity
        err = test_query_user_signin_activity()
    case TEST_CMD_CODE_PPROF_QUERY_SIGNIN_ACTIVITY: //pprof query user signin activity
        err = test_pprof_query_user_signin_activity(&done, &errSum, id, sum)
    case TEST_CMD_CODE_SIGNIN_ACTIVITY: //signin
        err = test_user_signin()
    case TEST_CMD_CODE_PPROF_SIGNIN_ACTIVITY: //pprof signin
        err = test_pprof_user_signin(&done, &errSum, id, sum)
    //checkpoint
    case TEST_CMD_CODE_QUERY_USERCHECKPOINT: // query user checkpoints
        err = test_query_checkpoints()
    case TEST_CMD_CODE_PPROF_QUERY_USERCHECKPOINT: // pprof query user checkpoints
        err = test_pprof_query_checkpoints(&done, &errSum, id, sum)
    case TEST_CMD_CODE_QUERY_CHECKPOINT_FRINENDRANK: //query user checkpoint friend rank
        err = test_query_checkpoint_friend_rank()
    case TEST_CMD_CODE_PPROF_QUERY_CHECKPOINT_FRINENDRANK: //pprof query user checkpoint friend rank
        err = test_pprof_query_checkpoint_friend_rank(&done, &errSum, id, sum)
    case TEST_CMD_CODE_QUERY_CHECKPOINT_GLOBALRANK: //query user checkpoint global rank
        err = test_query_checkpoint_global_rank()
    case TEST_CMD_CODE_PPROF_QUERY_CHECKPOINT_GLOBALRANK: //query user checkpoint global rank
        err = test_pprof_query_checkpoint_global_rank(&done, &errSum, id, sum)
    //wallet
    case TEST_CMD_CODE_QUERY_WALLET: // query user wallet
        err = test_query_user_wallet()
    case TEST_CMD_CODE_PPROF_QUERY_WALLET: // pprof query user wallet
        err = test_pprof_query_user_wallet(&done, &errSum, id, sum)
    //beforegame goods
    case TEST_CMD_CODE_QUERY_BEFOREGAMEGOODS: // query beforegame goods
        err = test_query_beforegame_goods()
    case TEST_CMD_CODE_PPROF_QUERY_BEFOREGAMEGOODS: // pprof query beforegame goods
        err = test_pprof_query_beforegame_goods(&done, &errSum, id, sum)
    case TEST_CMD_CODE_BUY_BEFOREGAMEGOODS: // buy beforegame goods
        err = test_buy_beforegame_goods()
    case TEST_CMD_CODE_PPROF_BUY_BEFOREGAMEGOODS: // pprof buy beforegame goods
        err = test_pprof_buy_beforegame_goods(&done, &errSum, id, sum)
    case TEST_CMD_CODE_USE_BEFOREGAMEGOODS: // use beforegame goods
        err = test_use_beforegame_goods()
    case TEST_CMD_CODE_PPROF_USE_BEFOREGAMEGOODS: // pprof use beforegame goods
        //todo
    //role
    case TEST_CMD_CODE_QUERY_USER_ROLEINFO: // query user roleinfo
        err = test_query_user_role_info()
    case TEST_CMD_CODE_PPROF_QUERY_USER_ROLEINFO: // pprof query user roleinfo
        err = test_pprof_user_role_info(&done, &errSum, id, sum)
    case TEST_CMD_CODE_QUERY_FRIEND_ROLEINFO: // query friend roleinfo
        err = test_query_friend_role_info()
    case TEST_CMD_CODE_PPROF_QUERY_FRIEND_ROLEINFO: // pprof query friend roleinfo
        err = test_query_friend_role_info()
    case TEST_CMD_CODE_SET_USER_SELECTED_ROLEINFO: // set user select role
        err = test_set_user_select_role_id()
    case TEST_CMD_CODE_PPROF_SET_USER_SELECTED_ROLEINFO: // pprof set user select role
        //todo
    case TEST_CMD_CODE_BUY_ROLE: // upgrade/buy user role
        err = test_upgrade_user_role()
    case TEST_CMD_CODE_PPROF_BUY_ROLE: // pprof upgrade/buy user role
        //todo
    //jigsaw
    case TEST_CMD_CODE_QUERY_USER_JIGSAW: // query user jigsaw
        err = test_query_user_jigsaw()
    case TEST_CMD_CODE_PPROF_QUERY_USER_JIGSAW: // pprof query user jigsaw
        //todo
    case TEST_CMD_CODE_BUY_JIGSAW: // buy user jigsaw
        err = test_buy_user_jigsaw()
    case TEST_CMD_CODE_PPROF_BUY_JIGSAW: // pprof buy user jigsaw
        //todo
    //friendmail
    case TEST_CMD_CODE_QUERY_FRIENDMAIL: // query user friend mails
        err = test_query_user_friend_mails()
    case TEST_CMD_CODE_PPROF_QUERY_FRIENDMAIL: // pprof query user friend mails
        err = test_pprof_query_friend_from_friendship(&done, &errSum, id, sum)
    case TEST_CMD_CODE_GIVE_FRIEND_FROM_FRIENDSHIP: // give friend from friendship
        err = test_give_friend_from_friendship()
    case TEST_CMD_CODE_PPROF_GIVE_FRIEND_FROM_FRIENDSHIP: // pprof give friend from friendship
        err = test_pprof_give_friend_from_friendship(&done, &errSum, id, sum)
    case TEST_CMD_CODE_APPLY_FRIEND_FROM_FRIENDSHIP: // apply friend from friendship
        err = test_apply_friend_from_friendship()
    case TEST_CMD_CODE_PPROF_APPLY_FRIEND_FROM_FRIENDSHIP: // pprof give friend from friendship
        err = test_pprof_apply_friend_from_friendship(&done, &errSum, id, sum)
    case TEST_CMD_CODE_COLLECT_FRIEND_FROM_FRIENDMAIL: // test confirm user friendmails
        err = test_confirm_user_friendmails()
    case TEST_CMD_CODE_PPROF_COLLECT_FRIEND_FROM_FRIENDMAIL: // pprof test confirm user friendmails
        err = test_pprof_confirm_friend_from_friendship(&done, &errSum, id, sum)
    //sysmail
    case TEST_CMD_CODE_QUERY_USER_SYSMAIL: // query user sysmails
        err = test_query_user_sysmail()
    case TEST_CMD_CODE_PPROF_QUERY_USER_SYSMAIL: // pprof query user sysmails
        err = test_pprof_query_user_sysmail(&done, &errSum, id, sum)
    case TEST_CMD_CODE_CONFIRM_USER_SYSMAIL: // confirm user sysmails
        err = test_confirm_user_sysmail()
    case TEST_CMD_CODE_PPROF_CONFIRM_USER_SYSMAIL: // pprof confirm user sysmails
        err = test_pprof_confirm_user_sysmail(&done, &errSum, id, sum)
    case TEST_CMD_CODE_READ_USER_SYSMAIL: // read user sysmails
        err = test_read_user_sysmail()
    case TEST_CMD_CODE_PPROF_READ_USER_SYSMAIL:
        err = test_pprof_read_user_sysmail(&done, &errSum, id, sum)
    case TEST_CMD_CODE_QUERY_RUNE: // query user runes
        err = test_query_user_runes()
    case TEST_CMD_CODE_PPROF_QUERY_RUNE: // pprof query user runes
    //todo
    case TEST_CMD_CODE_BUY_RUNE: // buy user rune
        err = test_buy_user_rune()
    case TEST_CMD_CODE_MEMCACHE_GET:
        err = test_memcache_get()
    case TEST_CMD_CODE_PPROF_MEMCACHE_GET:
        err = test_pprof_memcache_get(&done, &errSum, id, sum)
    case TEST_CMD_CODE_MEMCACHE_SET:
        err = test_memcache_set()
    case TEST_CMD_CODE_PPROF_MEMCACHE_SET:
        err = test_pprof_memcache_set(&done, &errSum, id, sum)
    case TEST_CMD_CODE_MEMCACHES_GET:
        err = test_memcaches_get()
    case TEST_CMD_CODE_MEMCACHES_SET:
        err = test_memcaches_set()
    case TEST_CMD_CODE_DB_QUERY_USER_ACCOUNT_SOMEFIELD:
        err = test_db_query_user_account_somefields()
    case TEST_CMD_CODE_MAINTENANCE_PROP:
        err = test_maintenance_prop()
    case TEST_CMD_CODE_APN_NOTIFICATION:
        err = test_apn_notification()
    case TEST_CMD_CODE_APN_NOTIFICATION2APNS:
        err = test_apn_notification2apns()
    case TEST_CMD_CODE_APN_ENABLEDEVICETOKEN2APNS:
        err = test_apn_enabledevicetoken2apns()
    case TEST_CMD_CODE_APN_DISABLEDEVICETOKEN2APNS:
        err = test_apn_disabledevicetoken2apns()
    case TEST_CMD_CODE_TIMER:
        err = test_timer()
    case TEST_CMD_CODE_DB:
        err = test_db_query_user_account_distinct()
    case TEST_CMD_CODE_ADVERTISEMENT:
        err = test_advertisement()
    case TEST_CMD_CODE_IAPSTATISTIC:
        err = test_iapstatistic()
    default:
    }

    if err != nil {
        xylog.ErrorNoId("======== test failed =========")
    } else {
        xylog.DebugNoId("======== test succeed =========")
    }

    time.Sleep(time.Second * 5)

}
