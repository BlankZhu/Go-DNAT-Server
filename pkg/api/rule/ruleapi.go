package apirule

import (
	memstorage "BlankZhu/Go-DNAT-Server/pkg/db/memory"
	mysqlstorage "BlankZhu/Go-DNAT-Server/pkg/db/mysql"
	"BlankZhu/Go-DNAT-Server/pkg/entity"
	"BlankZhu/Go-DNAT-Server/pkg/route"
	"BlankZhu/Go-DNAT-Server/pkg/util"
	"fmt"
	"net/http"

	"github.com/coreos/go-iptables/iptables"
	"github.com/gin-gonic/gin"
)

var mem *memstorage.MemStorage
var mysql *mysqlstorage.MySQLStorage

func init() {
	mem = memstorage.Get()
	mysql = mysqlstorage.Get()
}

func GetByID(c *gin.Context) {
	ruleID := c.Params.ByName("RuleID")
	value, ok := mem.Get(ruleID)
	if ok {
		c.JSON(http.StatusOK, value)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, mem.List())
}

func AddOne(c *gin.Context) {
	rule := entity.Rule{}
	err := c.BindJSON(&rule)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	if len(rule.CIDR) == 0 || len(rule.Dest) == 0 {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	ruleID := util.GetUUID()
	rule.RuleID = ruleID
	// try adding to iptables
	route.AddOneWithProtocol(&rule, iptables.ProtocolIPv4) // FIXME: Support ipv6
	// add to memory
	mem.Add(&rule)
	// add to db
	err = mysql.Add(&rule)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusServiceUnavailable)
		return
	}
	c.JSON(http.StatusOK, rule)
}

func DeleteByID(c *gin.Context) {
	ruleID := c.Params.ByName("RuleID")
	r, ok := mem.Get(ruleID)
	if ok {
		// delete from the iptables
		route.DeleteOneWithProtocol(r, iptables.ProtocolIPv4) // FIXME: Support ipv6
		// delete from memory
		mem.Delete(ruleID)
		// delete from db
		err := mysql.Delete(ruleID)
		if err != nil {
			fmt.Println(err.Error())
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusNoContent)
		return
	}

	c.Status(http.StatusNotFound)
	return
}

func UpdateByID(c *gin.Context) {
	ruleID := c.Params.ByName("RuleID")
	rule := entity.Rule{}
	err := c.BindJSON(&rule)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	r, ok := mem.Get(ruleID)
	if ok {
		// update iptables, usually NAT destination
		route.UpdateOneWithProtocol(r, iptables.ProtocolIPv4) // FIXME: Support ipv6
		// update the memory
		mem.Update(&rule)
		// update the database
		rule.RuleID = ruleID
		mysql.Update(&rule)
		if err != nil {
			fmt.Println(err.Error())
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.JSON(http.StatusOK, rule)
		return
	}

	c.Status(http.StatusNotFound)
	return
}
