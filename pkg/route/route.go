package route

import (
	"BlankZhu/Go-DNAT-Server/pkg/config"
	"BlankZhu/Go-DNAT-Server/pkg/entity"
	"BlankZhu/Go-DNAT-Server/pkg/util"
	"fmt"
	"sync"

	"github.com/coreos/go-iptables/iptables"
)

// TODO: add all those iptables rules to self-defined chains

var mtx sync.Mutex

// AddOneWithProtocol add a group of rules related to NAT to iptables
func AddOneWithProtocol(rule *entity.Rule, proto iptables.Protocol) error {
	mtx.Lock()
	defer mtx.Unlock()
	return addOneWithProtocol(rule, proto)
}

func addOneWithProtocol(rule *entity.Rule, proto iptables.Protocol) error {
	ipt, err := iptables.NewWithProtocol(proto)
	if err != nil {
		return err
	}
	conf := config.Get()

	preRuleSpec := genPreRuleSpec(conf.PublicIP, rule.Dest)
	postRuleSpec := genPostRuleSpec(conf.PrivateIP, rule.Dest)
	outputRuleSpec := genOutputRuleSpec(conf.PublicIP, rule.Dest)
	inputRuleSpec := genInputRuleSpec(rule.CIDR)

	err = ipt.AppendUnique(util.TableNAT, util.PreroutingChain, preRuleSpec)
	if err != nil {
		return err
	}
	err = ipt.AppendUnique(util.TableNAT, util.PostroutingChain, postRuleSpec)
	if err != nil {
		return err
	}
	err = ipt.AppendUnique(util.TableNAT, util.OutputChain, outputRuleSpec)
	if err != nil {
		return err
	}
	return ipt.AppendUnique(util.TableNAT, util.InputChain, inputRuleSpec)
}

// DeleteOneWithProtocol delete a group of rules related to NAT from iptables
func DeleteOneWithProtocol(rule *entity.Rule, proto iptables.Protocol) error {
	mtx.Lock()
	defer mtx.Unlock()
	return deleteOneWithProtocol(rule, proto)
}

func deleteOneWithProtocol(rule *entity.Rule, proto iptables.Protocol) error {
	ipt, err := iptables.NewWithProtocol(proto)
	if err != nil {
		return err
	}
	conf := config.Get()

	preRuleSpec := genPreRuleSpec(conf.PublicIP, rule.Dest)
	postRuleSpec := genPostRuleSpec(conf.PrivateIP, rule.Dest)
	outputRuleSpec := genOutputRuleSpec(conf.PublicIP, rule.Dest)
	inputRuleSpec := genInputRuleSpec(rule.CIDR)

	err = ipt.Delete(util.TableNAT, util.PreroutingChain, preRuleSpec)
	if err != nil {
		return err
	}
	err = ipt.Delete(util.TableNAT, util.PostroutingChain, postRuleSpec)
	if err != nil {
		return err
	}
	err = ipt.Delete(util.TableNAT, util.OutputChain, outputRuleSpec)
	if err != nil {
		return err
	}
	return ipt.Delete(util.TableNAT, util.InputChain, inputRuleSpec)
}

// UpdateOneWithProtocol combine AddOneWithProtocol and DeleteOneWithProtocol
func UpdateOneWithProtocol(rule *entity.Rule, proto iptables.Protocol) error {
	mtx.Lock()
	defer mtx.Unlock()
	err := deleteOneWithProtocol(rule, proto)
	if err != nil {
		return err
	}
	return addOneWithProtocol(rule, proto)
}

func genPreRuleSpec(publicIP, serverIP string) string {
	return fmt.Sprintf("--dst %s -j DNAT --to-destination %s", publicIP, serverIP)
}

func genPostRuleSpec(privateIP, serverIP string) string {
	return fmt.Sprintf("--dst %s -j SNAT --to-source %s", serverIP, privateIP)
}

func genOutputRuleSpec(publicIP, serverIP string) string {
	return fmt.Sprintf("--dst %s -j DNAT --to-destination %s", publicIP, serverIP)
}

func genInputRuleSpec(cidr string) string {
	return fmt.Sprintf("-s %s -j ACCEPT", cidr)
}
