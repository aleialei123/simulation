package fabric

type FabricConfig struct {
	Enabled                       bool   `mapstructure:"enabled"`
	OrderGeneralListenStartPort   int    `mapstructure:"order_general_listen_start_port"`
	OrderAdminListenStartPort     int    `mapstructure:"order_admin_listen_start_port"`
	OrderOperationListenStartPort int    `mapstructure:"order_operation_listen_start_port"`
	OrderStartPort                int    `mapstructure:"order_start_port"`
	OrderP2pStartPort			  int    `mapstructure:"order_p2p_start_port"`
	OrderRpcStartPort			  int    `mapstructure:"order_rpc_start_port"`
	PeerListenStartPort           int    `mapstructure:"peer_listen_start_port"`
	PeerChaincodeStartPort        int    `mapstructure:"peer_chaincode_start_port"`
	PeerOperationStartPort        int    `mapstructure:"peer_operation_start_port"`
	PeerP2pStartPort			  int    `mapstructure:"peer_p2p_start_port"`
	PeerRpcStartPort			  int    `mapstructure:"peer_rpc_start_port"`
	ConsensusType                 int    `mapstructure:"consensus_type"`
	LogLevel                      string `mapstructure:"log_level"`
	EnableBroadcastDefence        bool   `mapstructure:"enable_broadcast_defence"`
	DirectRemoveAttackedNode      bool   `mapstructure:"direct_remove_attacked_node"`
	SpeedCheck                    bool   `mapstructure:"speed_check"`
	FabricProjectPath             string `mapstructure:"fabric_project_path"`
	FabricNetworkPath             string `mapstructure:"fabric_network_path"`
}
