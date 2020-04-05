package cloudyipam

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Database string
	Port     string
	TLS      bool
}

// DSN assembles the various parts of a DatabaseConfig object into a
// PostgreSQL DSN.
func (d DatabaseConfig) DSN() string {
	dsn := []string{}
	if d.User != "" {
		dsn = append(dsn, fmt.Sprintf("user=%s", d.User))
	}
	if d.Password != "" {
		dsn = append(dsn, fmt.Sprintf("password='%s'", d.Password))
	}
	if d.Host != "" {
		dsn = append(dsn, fmt.Sprintf("host=%s", d.Host))
	}
	if d.Port != "" {
		dsn = append(dsn, fmt.Sprintf("port=%s", d.Port))
	}
	if d.Database != "" {
		dsn = append(dsn, fmt.Sprintf("dbname=%s", d.Database))
	}
	if d.TLS {
		dsn = append(dsn, "sslmode=verify-full")
	} else {
		dsn = append(dsn, "sslmode=disable")
	}
	return strings.Join(dsn, " ")
}

// CloudyIPAM client structure. No user-serviceable parts inside.
type CloudyIPAM struct {
	dsn  string
	conn *sql.DB
}

type Zone struct {
	Id        string
	Name      string
	Range     string
	PrefixLen int
}

type Subnet struct {
	Id        string
	Name      string
	Zone      string
	Range     string
	Available string
	Usage     string
}

// NewCloudyIPAM creates a CloudyIPAM client structure and connects to the
// specified PostgreSQL database
func NewCloudyIPAM(dsn string) (*CloudyIPAM, error) {
	obj := &CloudyIPAM{dsn: dsn}
	var err error
	obj.conn, err = sql.Open("postgres", obj.dsn)
	return obj, err
}

func (ipam CloudyIPAM) Ping() error {
	return ipam.conn.Ping()
}

// CreateZone initializes a new IPAM zone and returns its unique identifier
func (ipam CloudyIPAM) CreateZone(z Zone) (string, error) {
	var ident string
	err := ipam.conn.QueryRow("SELECT * FROM create_zone($1,$2,$3)", z.Name, z.Range, z.PrefixLen).Scan(&ident)
	if err != nil {
		return "", fmt.Errorf("Unable to create zone %s as %s subdivided into /%d subnets: %v", z.Name, z.Range, z.PrefixLen, err)
	}
	return ident, nil
}

type ReadZoneNotFoundError struct {
	Id string
}

func (e ReadZoneNotFoundError) Error() string {
	return fmt.Sprintf("Zone not found in database: %v", e.Id)
}

type ReadSubnetNotFoundError struct {
	Id string
}

func (e ReadSubnetNotFoundError) Error() string {
	return fmt.Sprintf("Subnet %s not found in database", e.Id)
}

type ZoneFullError struct {
	Id string
}

func (e *ZoneFullError) Error() string {
	return fmt.Sprintf("No available subnets remaining in zone %v", e.Id)
}

// ReadZone attempts to retrieve a named IPAM zone record. Mainly for Terraform's
// usage, not really that useful.
func (ipam CloudyIPAM) ReadZone(id string) (*Zone, error) {
	var zone, cidr string
	var prefixlen int
	err := ipam.conn.QueryRow("SELECT name,range,prefixlen FROM read_zone($1)", id).Scan(&zone, &cidr, &prefixlen)
	switch err {
	case nil: // found a zone, return it
		return &Zone{Id: id, Name: zone, Range: cidr, PrefixLen: prefixlen}, nil
	case sql.ErrNoRows: // query ok, no zone found
		return nil, &ReadZoneNotFoundError{Id: id}
	default: // some other error
		return nil, err
	}
}

// ListZones attempts to retrieve all zones.
func (ipam CloudyIPAM) ListZones() ([]Zone, error) {
	rows, err := ipam.conn.Query("SELECT * FROM list_zones()")
	if err != nil {
		return nil, fmt.Errorf("unable to list zones: %v", err)
	}
	zones := []Zone{}
	defer rows.Close()
	for rows.Next() {
		z := Zone{}
		if err := rows.Scan(&z.Id, &z.Name, &z.Range, &z.PrefixLen); err != nil {
			return nil, fmt.Errorf("unable to scan zones: %v", err)
		}
		zones = append(zones, z)
	}
	return zones, nil
}

// PopulateZone fills an empty IPAM zone with equal-sized allocatable subnets.
// Should not normally be required as this is invoked automatically when a zone
// is created.
func (ipam CloudyIPAM) PopulateZone(zone string) error {
	_, err := ipam.conn.Exec("CALL populate_zone($1)", zone)
	return err
}

// DestroyZone attempts to delete an IPAM zone. This will fail unless all
// subnets in the zone have been freed (with #FreeSubnet) and destroyed
// (with #DestroySubnet)
func (ipam CloudyIPAM) DestroyZone(zone string) error {
	_, err := ipam.conn.Exec("CALL destroy_zone($1)", zone)
	return err
}

// AllocateSubnet attempts to allocate a subnet in an IPAM zone. If none are
// available, an error is returned.
func (ipam CloudyIPAM) AllocateSubnet(zone string, usage string) (*Subnet, error) {
	var s Subnet
	err := ipam.conn.QueryRow("SELECT * FROM allocate_subnet($1,$2)", zone, usage).Scan(&s.Id, &s.Zone, &s.Range, &s.Available, &s.Usage)
	switch err {
	case nil: // allocated a subnet, return it
		return &s, nil
	case sql.ErrNoRows: // query executed but no available subnets
		return nil, &ZoneFullError{Id: zone}
	}
	return nil, err // some other error
}

// ReadSubnet attempts to read an allocated subnet in an IPAM zone. If no
// matching subnet is found, an error is returned.
func (ipam CloudyIPAM) ReadSubnet(subnet string) (*Subnet, error) {
	var s Subnet
	err := ipam.conn.QueryRow("SELECT * FROM read_subnet($1)", subnet).Scan(&s.Id, &s.Zone, &s.Range, &s.Available, &s.Usage)
	switch err {
	case nil: // found a subnet, return it
		return &s, nil
	case sql.ErrNoRows: // query ok, no zone found
		return nil, &ReadSubnetNotFoundError{Id: subnet}
	}
	return nil, err
}

// ListSubnets attempts to retrieve all subnets.
func (ipam CloudyIPAM) ListSubnets() ([]Subnet, error) {
	rows, err := ipam.conn.Query("SELECT * FROM list_subnets()")
	if err != nil {
		return nil, fmt.Errorf("unable to list subnets: %v", err)
	}
	subnets := []Subnet{}
	defer rows.Close()
	for rows.Next() {
		s := Subnet{}
		if err := rows.Scan(&s.Id, &s.Zone, &s.Range, &s.Available, &s.Usage); err != nil {
			return nil, fmt.Errorf("unable to scan subnets: %v", err)
		}
		subnets = append(subnets, s)
	}
	return subnets, nil
}

// FreeSubnet attempts to mark a subnet in an IPAM zone as allocatable.
func (ipam CloudyIPAM) FreeSubnet(subnet string) error {
	_, err := ipam.conn.Exec("CALL free_subnet($1)", subnet)
	return err
}

// DestroySubnet attempts to destroy a subnet. This will fail if the subnet
// is currently allocated.
func (ipam CloudyIPAM) DestroySubnet(subnet string) error {
	_, err := ipam.conn.Exec("CALL destroy_subnet($1)", subnet)
	return err
}
