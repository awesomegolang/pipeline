// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package defaults

import (
	"fmt"
	"time"

	"github.com/banzaicloud/pipeline/config"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	pkgErrors "github.com/banzaicloud/pipeline/pkg/errors"
	oracle "github.com/banzaicloud/pipeline/pkg/providers/oracle/model"
	"github.com/spf13/viper"
)

// cluster profile table names
const (
	DefaultEC2ProfileTableName         = "profiles_ec2"
	DefaultEC2NodePoolProfileTableName = "profiles_nodepools_ec2"

	DefaultEKSProfileTableName         = "profiles_eks"
	DefaultEKSNodePoolProfileTableName = "profiles_nodepools_eks"

	DefaultAKSProfileTableName         = "profiles_aks"
	DefaultAKSNodePoolProfileTableName = "profiles_nodepools_aks"

	DefaultGKEProfileTableName         = "profiles_gke"
	DefaultGKENodePoolProfileTableName = "profiles_nodepools_gke"
)

// default node name for all provider
const (
	DefaultNodeName = "pool1"
)

// SetDefaultValues saves the default cluster profile into the database if not exists yet
func SetDefaultValues() error {
	log.Info("setting up default cluster profiles")

	defaults := GetDefaultProfiles()

	for _, d := range defaults {
		if !d.IsDefinedBefore() { // the table not contains the default profile
			log.WithField("cloud", d.GetCloud()).Info("default profile is missing. Setting up...")

			if err := d.SaveInstance(); err != nil {
				return fmt.Errorf("could not save default values[%s]: %s", d.GetCloud(), err.Error())
			}
		} else { // default profile already exists
			log.WithField("cloud", d.GetCloud()).Info("default profile is already set up")
		}
	}

	return nil
}

// ClusterProfile describes a cluster profile
type ClusterProfile interface {
	IsDefinedBefore() bool
	SaveInstance() error
	GetCloud() string
	GetDistribution() string
	GetProfile() *pkgCluster.ClusterProfileResponse
	UpdateProfile(*pkgCluster.ClusterProfileRequest, bool) error
	DeleteProfile() error
}

// DefaultModel describes the common variables all types of clouds
type DefaultModel struct {
	Name      string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// save saves the given data into database
func save(i interface{}) error {
	database := config.DB()
	return database.Save(i).Error
}

// GetDefaultProfiles returns all types of clouds with default profile name.
func GetDefaultProfiles() []ClusterProfile {
	return []ClusterProfile{
		&EC2Profile{
			DefaultModel: DefaultModel{Name: GetDefaultProfileName()},
			NodePools: []*EC2NodePoolProfile{{
				AmazonNodePoolProfileBaseFields: AmazonNodePoolProfileBaseFields{
					NodeName: DefaultNodeName,
				},
			}},
		},
		&EKSProfile{
			DefaultModel: DefaultModel{Name: GetDefaultProfileName()},
			NodePools: []*EKSNodePoolProfile{{
				AmazonNodePoolProfileBaseFields: AmazonNodePoolProfileBaseFields{
					NodeName: DefaultNodeName,
				},
			}},
		},
		&AKSProfile{
			DefaultModel: DefaultModel{Name: GetDefaultProfileName()},
			NodePools: []*AKSNodePoolProfile{{
				NodeName: DefaultNodeName,
			}},
		},
		&GKEProfile{
			DefaultModel: DefaultModel{Name: GetDefaultProfileName()},
			NodePools: []*GKENodePoolProfile{{
				NodeName: DefaultNodeName,
			}},
		},
		&oracle.Profile{
			Name: GetDefaultProfileName(),
			NodePools: []*oracle.ProfileNodePool{{
				Name: DefaultNodeName,
			}},
		},
	}
}

// GetAllProfiles loads all saved cluster profile from database by given cloud type
func GetAllProfiles(distribution string) ([]ClusterProfile, error) {

	var defaults []ClusterProfile
	db := config.DB()

	switch distribution {

	case pkgCluster.EC2:
		var awsProfiles []EC2Profile
		db.Find(&awsProfiles)
		for i := range awsProfiles {
			defaults = append(defaults, &awsProfiles[i])
		}

	case pkgCluster.EKS:
		var eksProfiles []EKSProfile
		db.Find(&eksProfiles)
		for i := range eksProfiles {
			defaults = append(defaults, &eksProfiles[i])
		}

	case pkgCluster.AKS:
		var aksProfiles []AKSProfile
		db.Find(&aksProfiles)
		for i := range aksProfiles {
			defaults = append(defaults, &aksProfiles[i])
		}

	case pkgCluster.GKE:
		var gkeProfiles []GKEProfile
		db.Find(&gkeProfiles)
		for i := range gkeProfiles {
			defaults = append(defaults, &gkeProfiles[i])
		}

	case pkgCluster.OKE:
		okeProfiles := oracle.GetProfiles()
		for i := range okeProfiles {
			defaults = append(defaults, &okeProfiles[i])
		}

	default:
		return nil, pkgErrors.ErrorNotSupportedCloudType
	}

	return defaults, nil

}

// GetProfile finds cluster profile from database by given name and cloud type
func GetProfile(distribution string, name string) (ClusterProfile, error) {
	db := config.DB()

	switch distribution {
	case pkgCluster.EC2:
		var ec2Profile EC2Profile
		if err := db.Where(EC2Profile{DefaultModel: DefaultModel{Name: name}}).First(&ec2Profile).Error; err != nil {
			return nil, err
		}
		return &ec2Profile, nil

	case pkgCluster.EKS:
		var ec2Profile EKSProfile
		if err := db.Where(EKSProfile{DefaultModel: DefaultModel{Name: name}}).First(&ec2Profile).Error; err != nil {
			return nil, err
		}
		return &ec2Profile, nil

	case pkgCluster.AKS:
		var aksProfile AKSProfile
		if err := db.Where(GKEProfile{DefaultModel: DefaultModel{Name: name}}).First(&aksProfile).Error; err != nil {
			return nil, err
		}
		return &aksProfile, nil

	case pkgCluster.GKE:
		var gkeProfile GKEProfile
		if err := db.Where(GKEProfile{DefaultModel: DefaultModel{Name: name}}).First(&gkeProfile).Error; err != nil {
			return nil, err
		}
		return &gkeProfile, nil

	case pkgCluster.OKE:
		var okeProfile oracle.Profile
		okeProfile, err := oracle.GetProfileByName(name)
		return &okeProfile, err

	default:
		return nil, pkgErrors.ErrorNotSupportedCloudType
	}

}

// GetDefaultProfileName reads the default profile name env var
func GetDefaultProfileName() string {
	return viper.GetString("cloud.defaultProfileName")
}
