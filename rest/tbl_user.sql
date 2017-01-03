create sequence xx_user_seq increment by 1 minvalue 1 no maxvalue start with 1;
CREATE TABLE "xx_user" (
"id" int4 DEFAULT nextval('user_seq'::regclass) NOT NULL,
"user_code"  varchar(10) NOT NULL,
"user_name" varchar(30) NOT NULL,
"pwd" varchar(20) NOT NULL,
"update_time" timestamp(6) DEFAULT now() NOT NULL,
"create_time" timestamp(6) DEFAULT now() NOT NULL,
CONSTRAINT "xx_user_pkey" PRIMARY KEY ("id")
)
WITH (OIDS=FALSE)
;
ALTER TABLE "xx_user" OWNER TO "avcp_work";
COMMENT ON TABLE "xx_user" IS '用户表';
COMMENT ON COLUMN "xx_user"."user_code" IS '用户code';
COMMENT ON COLUMN "xx_user"."user_name" IS '用户名称';
COMMENT ON COLUMN "xx_user"."pwd" IS '用户密码';
CREATE INDEX "idx_xx_user_1" ON "xx_user" USING btree (user_code);

insert into xx_user(user_code,user_name,pwd) values('001','老李','001');
