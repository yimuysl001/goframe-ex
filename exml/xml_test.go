package exml

import (
	"fmt"
	"testing"
)

func TestXml(t *testing.T) {
	var cxml = `
<ClinicalDocument xmlns="urn:hl7-org:v3" xmlns:mif="urn:hl7-org:v3/mif" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:hl7-org:v3 ../../../core/sdschemas/SDA.xsd">
	<realmCode code="CN"/>
	<typeId root="2.16.840.1.113883.1.3" extension="POCD_MT000040"/>
	<templateId root="2.16.156.10011.2.1.1.22"/>
	<!--  文档流水号  -->
	<id root="2.16.156.10011.1.1" extension="1895655433323425793"/>
	<!--  code 默认值：C0002   （必填） -->
	<code code="C0002" codeSystem="2.16.156.10011.2.4" codeSystemName="卫生信息共享文档编码体系"/>
	<title>门（急）诊病历</title>
	<!--  文档机器生成时间 yyyyMMDDhhmmss（格式标椎）（必填） 用于签名串  -->
	<effectiveTime value="20250819112433"/>
	<confidentialityCode code="N" codeSystem="2.16.840.1.113883.5.25" codeSystemName="Confidentiality" displayName="正常访问保密级别"/>
	<languageCode code="zh-CN"/>
	<setId/>
	<versionNumber/>
	<recordTarget typeCode="RCT" contextControlCode="OP">
		<patientRole classCode="PAT">
			<!-- 门（急）诊号 （必填）院内唯一标识码  -->
			<id root="2.16.156.10011.1.11" extension="R202503000102"/>
			<patient classCode="PSN" determinerCode="INSTANCE">
				<!-- 患者身份证号 -->
				<id root="2.16.156.10011.1.3" extension="51322320100719221X"/>
				<!--  患者姓名 （必填）  -->
				<name>彭超</name>
				<!-- 患者性别 1男 2女 9未知 （必填） -->
				<administrativeGenderCode code="9" codeSystem="2.16.156.10011.2.3.3.4" codeSystemName="生理性别代码表（GB/T 2261.1）"/>
				<!--  患者出生年月 （必填）yyyyMMdd（标椎格式）没有为空字符传 -->
				<birthTime value="20100719"/>
				<!--  患者年龄 （必填）没有填0  -->
				<age unit="岁" value="14"/>
			</patient>
			<providerOrganization>
				<!--  就诊科室编号（必填）  -->
				<id root="2.16.156.10011.1.26" extension="A17"/>
				<!--  增加  @extension  就诊科室编号  -->
				<!--  就诊科室名称 （必填）  -->
				<name>结核病科</name>
				<asOrganizationPartOf>
					<wholeOrganization>
						<!--  机构代码（必填）来源于统计直报的机构编码（必须一致）  -->
						<id root="2.16.156.10011.1.5" extension="510000003696"/>
						<!--  机构名称 （必填）  -->
						<name>茂县人民医院</name>
					</wholeOrganization>
				</asOrganizationPartOf>
			</providerOrganization>
		</patientRole>
	</recordTarget>
	<!-- 创建者 -->
	<author typeCode="AUT" contextControlCode="OP">
		<!--  就诊日期时间 yyyyMMDDhhmmss（格式标椎）（必填）  -->
		<time value="20250301100041"/>
		<assignedAuthor classCode="ASSIGNED">
			<!--  医护人员身份证号（必填且格式必须是中国公民身份证格式）  -->
			<id root="2.16.156.10011.1.3" extension="513221198004150025"/>
			<!--  his系统医生编码 （必填）  -->
			<id root="2.16.156.10011.1.7" extension="1425722178066837506"/>
			<!--  就诊医师（必填）  -->
			<assignedPerson>
				<name>熊富丽</name>
			</assignedPerson>
		</assignedAuthor>
	</author>
	<!--  保管机构  -->
	<custodian typeCode="CST">
		<assignedCustodian classCode="ASSIGNED">
			<representedCustodianOrganization classCode="ORG" determinerCode="INSTANCE">
				<!--  机构代码（必填）来源于统计直报的机构编码（必须一致）  -->
				<id root="2.16.156.10011.1.5" extension="510000003696"/>
				<!--  机构名称 （必填）  -->
				<name>茂县人民医院</name>
			</representedCustodianOrganization>
		</assignedCustodian>
	</custodian>
	<!--  医师签名  -->
	<legalAuthenticator>
		<!--  就诊日期时间 yyyyMMDDhhmmss（格式标椎）（必填）  -->
		<time value="20250301100041"/>
		<!--  增加  @value属性  签名时间  -->
		<signatureCode/>
		<assignedEntity>
			<!--  医护人员身份证号（必填且格式必须是中国公民身份证格式）  -->
			<id root="2.16.156.10011.1.3" extension="513221198004150025"/>
			<!--  his系统医生编码 （必填）  -->
			<id root="2.16.156.10011.1.4" extension="1425722178066837506"/>
			<code displayName="责任医生"/>
			<!--  就诊医师（必填）  -->
			<assignedPerson>
				<name>熊富丽</name>
			</assignedPerson>
		</assignedEntity>
	</legalAuthenticator>
	<relatedDocument typeCode="RPLC">
		<parentDocument>
			<id/>
			<setId/>
			<versionNumber/>
		</parentDocument>
		<!--  增加 签名  -->
		<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
			<SignedInfo>
				<CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
				<SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
				<Reference>
					<DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
					<!--  数字签名字符串,没有可以为空  -->
					<DigestValue></DigestValue>
				</Reference>
			</SignedInfo>
			<!--  签名字符串MD5加密("ORGID={1}&APPID={2} &TIMESTAMP={3}&SN={4} 文档创建时间) （必填）  -->
			<SignatureValue>3ae3e7420fe926cd2baf4052df44f947</SignatureValue>
			<KeyInfo>
				<KeyValue>
					<DSAKeyValue>
						<!--  对称密钥，base64的字符串，没有传空  -->
						<Y></Y>
					</DSAKeyValue>
				</KeyValue>
			</KeyInfo>
		</Signature>
	</relatedDocument>
	<component>
		<structuredBody>
			<!--  过敏史章节  -->
			<component>
				<section>
					<code code="48765-2" displayName="Allergies, adverse reactions, alerts" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE02.10.023.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="过敏史标志"/>
							<value xsi:type="BL" value="false"/>
							<entryRelationship typeCode="COMP">
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE05.01.022.00" displayName="过敏史" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
									<value xsi:type="ST"></value>
								</observation>
							</entryRelationship>
						</observation>
					</entry>
				</section>
			</component>
			<!-- 主诉章节 -->
			<component>
				<section>
					<code code="10154-3" displayName="CHIEF COMPLAINT" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<!-- 主诉条目，有就传没有为空 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE04.01.119.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="主诉"/>
							<value xsi:type="ST">常规结核病筛查</value>
						</observation>
					</entry>
				</section>
			</component>
			<!-- 现病史章节 -->
			<component>
				<section>
					<code code="10164-2" displayName="HISTORY OF PRESENT ILLNESS" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<!-- 现病史条目，有就传没有为空 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE02.10.071.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="现病史"/>
							<value xsi:type="ST">无咳嗽、咳痰、潮热、盗汗、乏力、纳差、体重减轻等结核中毒症状。</value>
						</observation>
					</entry>
				</section>
			</component>
			<!--  既往史章节，有就传没有为空  -->
			<component>
				<section>
					<code code="11348-0" displayName="HISTORY OF PAST ILLNESS" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE02.10.099.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="既往史"/>
							<value xsi:type="ST">既往无结核病及其他病史。 </value>
						</observation>
					</entry>
				</section>
			</component>
			<!-- 体格检查章节，有就传没有为空 -->
			<component>
				<section>
					<code code="29545-1" displayName="PHYSICAL EXAMINATION" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<!-- 体格检查-一般状况检查结果 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE04.10.258.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="体格检查"/>
							<value xsi:type="ST">无特殊</value>
						</observation>
					</entry>
				</section>
			</component>
			<!--  实验室检验章节，有就传没有为空  -->
			<component>
				<section>
					<code code="30954-2" displayName="STUDIES SUMMARY" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<entry>
						<organizer classCode="CLUSTER" moodCode="EVN">
							<statusCode/>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE04.30.010.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="辅助检查项目"/>
								</observation>
							</component>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE04.30.009.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="辅助检查结果"/>
									<value xsi:type="ST"></value>
								</observation>
							</component>
						</organizer>
					</entry>
				</section>
			</component>
			<!--  诊断记录章节  -->
			<component>
				<section>
					<code code="29548-5" displayName="Diagnosis" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<!-- 初诊标志代码 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE06.00.196.00" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录" displayName="初诊标志代码"/>
							<value xsi:type="CD" code="1" codeSystem="2.16.156.10011.2.3.2.39" codeSystemName="初诊标志代码表" displayName="初诊"/>
						</observation>
					</entry>
					<!-- 中医“四诊”观察结果 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE02.10.028.00" displayName="中医“四诊”观察结果" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
							<value xsi:type="ST">常规结核病筛查</value>
						</observation>
					</entry>
					<!-- 条目：诊断 -->
							<entry>
								<organizer classCode="CLUSTER" moodCode="EVN">
									<statusCode/>
									<component>
										<observation classCode="OBS" moodCode="EVN">
											<!-- 主诊断名称（必填） -->
											<code code="DE05.01.025.00" displayName="诊断名称" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
											<value xsi:type="ST">健康查体</value>
										</observation>
									</component>
									<component>
										<observation classCode="OBS" moodCode="EVN">
											<!-- 主诊断代码（必填） -->
											<code code="DE05.01.024.00" displayName="诊断代码" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
											<value xsi:type="CD" code="Z00.001" codeSystem="2.16.156.10011.2.3.3.11.3" codeSystemName="诊断代码表（ICD-10）"/>
										</observation>
									</component>
								</organizer>
							</entry>
				</section>
			</component>
			<!--  治疗计划章节  -->
			<component>
				<section>
					<code code="18776-5" displayName="TREATMENT PLAN" codeSystem="2.16.840.1.113883.6.1" codeSystemName="LOINC"/>
					<text/>
					<!-- 辨证依据描述 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE05.10.132.00" displayName="辨证依据" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
							<value xsi:type="ST"></value>
						</observation>
					</entry>
					<!-- 治则治法 -->
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE06.00.300.00" displayName="治则治法" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
							<value xsi:type="ST"></value>
						</observation>
					</entry>
				</section>
			</component>
			<!-- 新增 卫生事件章节  -->
			<component>
				<section>
					<code displayName="卫生事件"/>
					<text/>
					<entry>
						<organizer classCode="BATTERY" moodCode="EVN">
							<statusCode/>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE09.00.250.00" displayName="医保类别" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
									<value xsi:type="CD" code="" displayName="fdsgd324532532"/>
								</observation>
							</component>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE07.00.007.00" displayName="参保类型" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
									<value xsi:type="CD" code="" displayName="fdadgdafdfad"/>
								</observation>
							</component>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE02.01.038.00" displayName="统筹区编码" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
									<value xsi:type="CD" code="" displayName=""/>
								</observation>
							</component>
						</organizer>
					</entry>
					<entry>
						<organizer classCode="BATTERY" moodCode="EVN">
							<statusCode/>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE02.10.014.00" displayName="怀孕标识" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
									<value xsi:type="BL" value="false"/>
								</observation>
							</component>
							<component>
								<observation classCode="OBS" moodCode="EVN">
									<code code="DE04.01.125.00" displayName="哺乳标识" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
									<value xsi:type="BL" value="false"/>
								</observation>
							</component>
						</organizer>
					</entry>
					<entry>
						<observation classCode="OBS" moodCode="EVN">
							<code code="DE05.00.158.00" displayName="特殊情况说明" codeSystem="2.16.156.10011.2.2.1" codeSystemName="卫生信息数据元目录"/>
							<value xsi:type="ST"></value>
						</observation>
					</entry>
				</section>
			</component>
		</structuredBody>
	</component>
</ClinicalDocument>
`

	root := GetElementRoot(cxml)

	value := GetElementValue(root, "./component/structuredBody/component[9]/section/entry[3]/observation/value/@xsi:type")
	fmt.Println(value)

}
