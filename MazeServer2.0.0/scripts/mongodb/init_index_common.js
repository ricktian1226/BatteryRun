// tpidmap ��ҵ������˻���Ϣ
db.tpidmap.ensureIndex({sid:1,source:1},{background:1});//���ݵ�����id��ѯ��Ҷ�Ӧtpid��Ϣ
db.tpidmap.ensureIndex({gid:1,source:1},{background:1});//����uid��ѯ��Ҷ�Ӧtpid��Ϣ
// useraccomplishment ��ҳɾ���Ϣ
db.useraccomplishment.ensureIndex({uid:1},{background:1});//��ʱ������ع���������Ϣ
// usercheckpoint ��Ҽ������Ϣ
db.usercheckpoint.ensureIndex({uid:1,checkpointid:1},{background:1});//��ѯ��Ҽ������Ϣ�б�
db.usercheckpoint.ensureIndex({checkpointid:1, score:-1},{background:1});//��ѯ��Ҽ������Ϣ�б�
// memcache ��һ�����Ϣ
db.memcache.ensureIndex({uid:1, key:1, platform:1},{background:1});
