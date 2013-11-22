// go run gen.go
// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package ipv4

// Internet Control Message Protocol (ICMP) Parameters, Updated: 2013-04-19
const (
	ICMPTypeEchoReply              ICMPType = 0  // Echo Reply
	ICMPTypeDestinationUnreachable ICMPType = 3  // Destination Unreachable
	ICMPTypeRedirect               ICMPType = 5  // Redirect
	ICMPTypeEcho                   ICMPType = 8  // Echo
	ICMPTypeRouterAdvertisement    ICMPType = 9  // Router Advertisement
	ICMPTypeRouterSolicitation     ICMPType = 10 // Router Solicitation
	ICMPTypeTimeExceeded           ICMPType = 11 // Time Exceeded
	ICMPTypeParameterProblem       ICMPType = 12 // Parameter Problem
	ICMPTypeTimestamp              ICMPType = 13 // Timestamp
	ICMPTypeTimestampReply         ICMPType = 14 // Timestamp Reply
	ICMPTypePhoturis               ICMPType = 40 // Photuris
)

// Internet Control Message Protocol (ICMP) Parameters, Updated: 2013-04-19
var icmpTypes = map[ICMPType]string{
	0:  "echo reply",
	3:  "destination unreachable",
	5:  "redirect",
	8:  "echo",
	9:  "router advertisement",
	10: "router solicitation",
	11: "time exceeded",
	12: "parameter problem",
	13: "timestamp",
	14: "timestamp reply",
	40: "photuris",
}

// Protocol Numbers, Updated: 2013-02-17
const (
	ianaProtocolIP             = 0   // IPv4 encapsulation, pseudo protocol number
	ianaProtocolHOPOPT         = 0   // IPv6 Hop-by-Hop Option
	ianaProtocolICMP           = 1   // Internet Control Message
	ianaProtocolIGMP           = 2   // Internet Group Management
	ianaProtocolGGP            = 3   // Gateway-to-Gateway
	ianaProtocolIPv4           = 4   // IPv4 encapsulation
	ianaProtocolST             = 5   // Stream
	ianaProtocolTCP            = 6   // Transmission Control
	ianaProtocolCBT            = 7   // CBT
	ianaProtocolEGP            = 8   // Exterior Gateway Protocol
	ianaProtocolIGP            = 9   // any private interior gateway (used by Cisco for their IGRP)
	ianaProtocolBBNRCCMON      = 10  // BBN RCC Monitoring
	ianaProtocolNVPII          = 11  // Network Voice Protocol
	ianaProtocolPUP            = 12  // PUP
	ianaProtocolARGUS          = 13  // ARGUS
	ianaProtocolEMCON          = 14  // EMCON
	ianaProtocolXNET           = 15  // Cross Net Debugger
	ianaProtocolCHAOS          = 16  // Chaos
	ianaProtocolUDP            = 17  // UserFields Datagram
	ianaProtocolMUX            = 18  // Multiplexing
	ianaProtocolDCNMEAS        = 19  // DCN Measurement Subsystems
	ianaProtocolHMP            = 20  // Host Monitoring
	ianaProtocolPRM            = 21  // Packet Radio Measurement
	ianaProtocolXNSIDP         = 22  // XEROX NS IDP
	ianaProtocolTRUNK1         = 23  // Trunk-1
	ianaProtocolTRUNK2         = 24  // Trunk-2
	ianaProtocolLEAF1          = 25  // Leaf-1
	ianaProtocolLEAF2          = 26  // Leaf-2
	ianaProtocolRDP            = 27  // Reliable Data Protocol
	ianaProtocolIRTP           = 28  // Internet Reliable Transaction
	ianaProtocolISOTP4         = 29  // ISO Transport Protocol Class 4
	ianaProtocolNETBLT         = 30  // Bulk Data Transfer Protocol
	ianaProtocolMFENSP         = 31  // MFE Network Services Protocol
	ianaProtocolMERITINP       = 32  // MERIT Internodal Protocol
	ianaProtocolDCCP           = 33  // Datagram Congestion Control Protocol
	ianaProtocol3PC            = 34  // Third Party Connect Protocol
	ianaProtocolIDPR           = 35  // Inter-Domain Policy Routing Protocol
	ianaProtocolXTP            = 36  // XTP
	ianaProtocolDDP            = 37  // Datagram Delivery Protocol
	ianaProtocolIDPRCMTP       = 38  // IDPR Control Message Transport Proto
	ianaProtocolTPPP           = 39  // TP++ Transport Protocol
	ianaProtocolIL             = 40  // IL Transport Protocol
	ianaProtocolIPv6           = 41  // IPv6 encapsulation
	ianaProtocolSDRP           = 42  // Source Demand Routing Protocol
	ianaProtocolIPv6Route      = 43  // Routing Header for IPv6
	ianaProtocolIPv6Frag       = 44  // Fragment Header for IPv6
	ianaProtocolIDRP           = 45  // Inter-Domain Routing Protocol
	ianaProtocolRSVP           = 46  // Reservation Protocol
	ianaProtocolGRE            = 47  // Generic Routing Encapsulation
	ianaProtocolDSR            = 48  // Dynamic Source Routing Protocol
	ianaProtocolBNA            = 49  // BNA
	ianaProtocolESP            = 50  // Encap Security Payload
	ianaProtocolAH             = 51  // Authentication Header
	ianaProtocolINLSP          = 52  // Integrated Net Layer Security  TUBA
	ianaProtocolSWIPE          = 53  // IP with Encryption
	ianaProtocolNARP           = 54  // NBMA Address Resolution Protocol
	ianaProtocolMOBILE         = 55  // IP Mobility
	ianaProtocolTLSP           = 56  // Transport Layer Security Protocol using Kryptonet key management
	ianaProtocolSKIP           = 57  // SKIP
	ianaProtocolIPv6ICMP       = 58  // ICMP for IPv6
	ianaProtocolIPv6NoNxt      = 59  // No Next Header for IPv6
	ianaProtocolIPv6Opts       = 60  // Destination Options for IPv6
	ianaProtocolCFTP           = 62  // CFTP
	ianaProtocolSATEXPAK       = 64  // SATNET and Backroom EXPAK
	ianaProtocolKRYPTOLAN      = 65  // Kryptolan
	ianaProtocolRVD            = 66  // MIT Remote Virtual Disk Protocol
	ianaProtocolIPPC           = 67  // Internet Pluribus Packet Core
	ianaProtocolSATMON         = 69  // SATNET Monitoring
	ianaProtocolVISA           = 70  // VISA Protocol
	ianaProtocolIPCV           = 71  // Internet Packet Core Utility
	ianaProtocolCPNX           = 72  // Computer Protocol Network Executive
	ianaProtocolCPHB           = 73  // Computer Protocol Heart Beat
	ianaProtocolWSN            = 74  // Wang Span Network
	ianaProtocolPVP            = 75  // Packet Video Protocol
	ianaProtocolBRSATMON       = 76  // Backroom SATNET Monitoring
	ianaProtocolSUNND          = 77  // SUN ND PROTOCOL-Temporary
	ianaProtocolWBMON          = 78  // WIDEBAND Monitoring
	ianaProtocolWBEXPAK        = 79  // WIDEBAND EXPAK
	ianaProtocolISOIP          = 80  // ISO Internet Protocol
	ianaProtocolVMTP           = 81  // VMTP
	ianaProtocolSECUREVMTP     = 82  // SECURE-VMTP
	ianaProtocolVINES          = 83  // VINES
	ianaProtocolTTP            = 84  // TTP
	ianaProtocolIPTM           = 84  // Protocol Internet Protocol Traffic Manager
	ianaProtocolNSFNETIGP      = 85  // NSFNET-IGP
	ianaProtocolDGP            = 86  // Dissimilar Gateway Protocol
	ianaProtocolTCF            = 87  // TCF
	ianaProtocolEIGRP          = 88  // EIGRP
	ianaProtocolOSPFIGP        = 89  // OSPFIGP
	ianaProtocolSpriteRPC      = 90  // Sprite RPC Protocol
	ianaProtocolLARP           = 91  // Locus Address Resolution Protocol
	ianaProtocolMTP            = 92  // Multicast Transport Protocol
	ianaProtocolAX25           = 93  // AX.25 Frames
	ianaProtocolIPIP           = 94  // IP-within-IP Encapsulation Protocol
	ianaProtocolMICP           = 95  // Mobile Internetworking Control Pro.
	ianaProtocolSCCSP          = 96  // Semaphore Communications Sec. Pro.
	ianaProtocolETHERIP        = 97  // Ethernet-within-IP Encapsulation
	ianaProtocolENCAP          = 98  // Encapsulation Header
	ianaProtocolGMTP           = 100 // GMTP
	ianaProtocolIFMP           = 101 // Ipsilon Flow Management Protocol
	ianaProtocolPNNI           = 102 // PNNI over IP
	ianaProtocolPIM            = 103 // Protocol Independent Multicast
	ianaProtocolARIS           = 104 // ARIS
	ianaProtocolSCPS           = 105 // SCPS
	ianaProtocolQNX            = 106 // QNX
	ianaProtocolAN             = 107 // Active Networks
	ianaProtocolIPComp         = 108 // IP Payload Compression Protocol
	ianaProtocolSNP            = 109 // Sitara Networks Protocol
	ianaProtocolCompaqPeer     = 110 // Compaq Peer Protocol
	ianaProtocolIPXinIP        = 111 // IPX in IP
	ianaProtocolVRRP           = 112 // Virtual Router Redundancy Protocol
	ianaProtocolPGM            = 113 // PGM Reliable Transport Protocol
	ianaProtocolL2TP           = 115 // Layer Two Tunneling Protocol
	ianaProtocolDDX            = 116 // D-II Data Exchange (DDX)
	ianaProtocolIATP           = 117 // Interactive Agent Transfer Protocol
	ianaProtocolSTP            = 118 // Schedule Transfer Protocol
	ianaProtocolSRP            = 119 // SpectraLink Radio Protocol
	ianaProtocolUTI            = 120 // UTI
	ianaProtocolSMP            = 121 // Simple Message Protocol
	ianaProtocolSM             = 122 // SM
	ianaProtocolPTP            = 123 // Performance Transparency Protocol
	ianaProtocolISIS           = 124 // ISIS over IPv4
	ianaProtocolFIRE           = 125 // FIRE
	ianaProtocolCRTP           = 126 // Combat Radio Transport Protocol
	ianaProtocolCRUDP          = 127 // Combat Radio UserFields Datagram
	ianaProtocolSSCOPMCE       = 128 // SSCOPMCE
	ianaProtocolIPLT           = 129 // IPLT
	ianaProtocolSPS            = 130 // Secure Packet Shield
	ianaProtocolPIPE           = 131 // Private IP Encapsulation within IP
	ianaProtocolSCTP           = 132 // Stream Control Transmission Protocol
	ianaProtocolFC             = 133 // Fibre Channel
	ianaProtocolRSVPE2EIGNORE  = 134 // RSVP-E2E-IGNORE
	ianaProtocolMobilityHeader = 135 // Mobility Header
	ianaProtocolUDPLite        = 136 // UDPLite
	ianaProtocolMPLSinIP       = 137 // MPLS-in-IP
	ianaProtocolMANET          = 138 // MANET Protocols
	ianaProtocolHIP            = 139 // Host Identity Protocol
	ianaProtocolShim6          = 140 // Shim6 Protocol
	ianaProtocolWESP           = 141 // Wrapped Encapsulating Security Payload
	ianaProtocolROHC           = 142 // Robust Header Compression
	ianaProtocolReserved       = 255 // Reserved
)
