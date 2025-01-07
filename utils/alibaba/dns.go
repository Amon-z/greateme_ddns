package alibaba

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

// GetAllDNSListByDomainNameAndRR 根据域名获取所有的DNS解析列表
func GetAllDNSListByDomainNameAndRR(domainName, rr *string) (*[]*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {
	// 最终结果
	var dnsList []*alidns20150109.DescribeDomainRecordsResponseBodyDomainRecordsRecord

	// 当前页码
	var currentPageNum int64 = 1

	for {
		// 初始化查询参数
		describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
			DomainName: domainName,
			PageNumber: &currentPageNum,
			PageSize:   tea.Int64(20),
			RRKeyWord:  rr,
		}
		runtime := &service.RuntimeOptions{}
		// 查询域名列表
		dnsResult, err := dnsClient.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, runtime)
		if err != nil {
			return nil, err
		}
		dnsList = append(dnsList, dnsResult.Body.DomainRecords.Record...)

		// 判断是否要继续分页
		if int(*dnsResult.Body.TotalCount) <= cap(dnsList) {
			// 查询到现在的域名数量大于等于总域名数量，那么我们就跳出循环
			break
		}
		// 继续循环，页码+1
		currentPageNum++
	}
	return &dnsList, nil
}

// AddDNSRecord 新增DNS解析记录
func AddDNSRecord(domain, rr, ipAddress *string, dnsType *string) error {
	// 封装修改的参数
	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: domain,
		RR:         rr,
		Type:       GetDNSType(dnsType),
		Value:      ipAddress,
		TTL:        tea.Int64(600),
		Line:       tea.String("default"),
	}
	runtime := &service.RuntimeOptions{}
	// 执行修改
	_, err := dnsClient.AddDomainRecordWithOptions(addDomainRecordRequest, runtime)
	return err
}

// UpdateDNSRecord 更新DNS解析记录
func UpdateDNSRecord(recordId, rr, ipAddress *string, dnsType *string) error {
	// 封装修改的参数
	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: recordId,
		RR:       rr,
		Type:     GetDNSType(dnsType),
		Value:    ipAddress,
		TTL:      tea.Int64(600),
		Line:     tea.String("default"),
	}
	runtime := &service.RuntimeOptions{}
	// 执行修改
	_, err := dnsClient.UpdateDomainRecordWithOptions(updateDomainRecordRequest, runtime)
	return err
}

// GetDNSType 获取解析的IP地址类型
func GetDNSType(dnsType *string) *string {
	recordType := tea.String("A")
	if strings.Compare("ipv6", *dnsType) == 0 {
		recordType = tea.String("AAAA")
	}
	return recordType
}
