/*
   mongodb�Ĺ̶�����������Ϊ��־���ݶ����ģ����е���־���ݿ⼯�϶����ù̶����Ϸ�ʽ������
   ��־���ϰ�����
   accountlog   ��¼��־
   gamelog      ��Ϸ��־
   shoppinglog  ������־
   iaplog       �ڹ���־
   lottolog     �齱��־
   maintenancelog ��Ӫ��־  
   collection��size��max��������ա�User����ģ��.xls���ĵ���
   https://192.168.1.204:8443/svn/Savage/Program/MazeServer2.0.0/scripts/mongodb
*/
db.createCollection("accountlog",{size:2147483648, capped:true, max:3000000,autoIndexId:false});
db.createCollection("gamelog",{size:2147483648, capped:true, max:5000000,autoIndexId:false});
db.createCollection("shoppinglog",{size:524288000, capped:true, max:3000000,autoIndexId:false});
db.createCollection("iaplog",{size:524288000, capped:true, max:3000000,autoIndexId:false});
db.createCollection("lottolog",{size:629145600, capped:true, max:3000000,autoIndexId:false});
db.createCollection("maintenancelog",{size:10485760, capped:true, max:3000,autoIndexId:false});
