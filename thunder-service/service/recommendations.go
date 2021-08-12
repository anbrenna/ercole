// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package service is a package that provides methods for querying data
package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/example/helpers"
	"github.com/oracle/oci-go-sdk/v45/optimizer"
)

// 0 high-availability-name high-availability-desc 0 ocid1.optimizercategory.oc1..aaaaaaaa34bnli7iarbwzpaz6abuup2egv2hcdyiixzellvxplhmceqvzula
// 1 performance-name performance-desc 0 ocid1.optimizercategory.oc1..aaaaaaaazihon4oggqfjsut576fsd45cm2d6g7g6qr2n3nxlpljsp7j33n4q
// 2 cost-management-name cost-management-desc 25.92 ocid1.optimizercategory.oc1..aaaaaaaaqeiskhuyp4pr7tohuooyujgyjmcq6cibc3btq6na62ev4ytz7ppa

func (as *ThunderService) GetOCRecommendations(compartmentId string) ([]model.Recommendation, error) {

	return GetOCListRecommendations("ocid1.tenancy.oc1..aaaaaaaazizzbqqbjv2se3y3fvm5osfumnorh32nznanirqoju3uks4buh4q")
}

func (as *ThunderService) GetOCRecommendationsWithCategory(compartmentId string) ([]model.RecommendationWithCategory, error) {

	return GetOCListRecommendationsWithCategory("ocid1.tenancy.oc1..aaaaaaaazizzbqqbjv2se3y3fvm5osfumnorh32nznanirqoju3uks4buh4q")
}

func GetOCListRecommendations(compartmentId string) ([]model.Recommendation, error) {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := optimizer.NewOptimizerClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	if err != nil {
		return nil, err
	} else {
		// Create a request and dependent object(s).
		req := optimizer.ListRecommendationsRequest{
			CompartmentId:          &compartmentId,
			CategoryId:             common.String("ocid1.optimizercategory.oc1..aaaaaaaaqeiskhuyp4pr7tohuooyujgyjmcq6cibc3btq6na62ev4ytz7ppa"),
			CompartmentIdInSubtree: common.Bool(true),
			Limit:                  common.Int(964),
		}

		// Send the request using the service client
		resp, err := client.ListRecommendations(context.Background(), req)
		helpers.FatalIfError(err)
		if err != nil {
			return nil, err
		} else {
			var cnt int
			var recTmp model.Recommendation
			var listRec []model.Recommendation

			// Retrieve value from the response.
			for i, s := range resp.Items {
				fmt.Println(i, *s.Name, *s.EstimatedCostSaving, s.Status, s.Importance, *s.Id)
				for j, p := range s.ResourceCounts {
					fmt.Println(j, *p.Count, p.Status)
					if p.Status == "PENDING" {
						cnt = *p.Count
					}
					recTmp = model.Recommendation{*s.Name, strconv.Itoa(cnt), fmt.Sprintf("%.2f", *s.EstimatedCostSaving), fmt.Sprintf("%v", s.Status), fmt.Sprintf("%v", s.Importance), *s.Id}

				}
				listRec = append(listRec, recTmp)
			}
			return listRec, nil
		}
	}
}

func GetOCListRecommendationsWithCategory(compartmentId string) ([]model.RecommendationWithCategory, error) {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := optimizer.NewOptimizerClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	if err != nil {
		return nil, err
	} else {
		var listRecWithCat []model.RecommendationWithCategory
		// Retrieve Recommendation Categories
		listCategory, err := GetOCListCategories(compartmentId)

		if err != nil {
			return nil, err
		} else {

			for _, q := range listCategory {
				// Create a request and dependent object(s).
				req := optimizer.ListRecommendationsRequest{
					CompartmentId:          &compartmentId,
					CategoryId:             &q.CategoryId,
					CompartmentIdInSubtree: common.Bool(true),
				}

				// Send the request using the service client
				resp, err := client.ListRecommendations(context.Background(), req)
				helpers.FatalIfError(err)
				if err != nil {
					return nil, err
				} else {
					var cnt int
					var recTmp model.Recommendation
					var recWithCatTmp model.RecommendationWithCategory
					var listRec []model.Recommendation

					// Retrieve value from the response.
					for _, s := range resp.Items {
						//fmt.Println(i, *s.Name, *s.EstimatedCostSaving, s.Status, s.Importance, *s.Id)
						for _, p := range s.ResourceCounts {
							//fmt.Println(j, *p.Count, p.Status)
							if p.Status == "PENDING" {
								cnt = *p.Count
							}
							recTmp = model.Recommendation{*s.Name, strconv.Itoa(cnt), fmt.Sprintf("%.2f", *s.EstimatedCostSaving), fmt.Sprintf("%v", s.Status), fmt.Sprintf("%v", s.Importance), *s.Id}

						}
						listRec = append(listRec, recTmp)
					}
					recWithCatTmp = model.RecommendationWithCategory{q.Name, listRec}
					listRecWithCat = append(listRecWithCat, recWithCatTmp)
				}
			}

			return listRecWithCat, nil
		}
	}
}

func GetOCListCategories(compartmentId string) ([]model.Category, error) {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := optimizer.NewOptimizerClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(err)

	// Create a request and dependent object(s).
	req := optimizer.ListCategoriesRequest{
		CompartmentId:          &compartmentId,
		CompartmentIdInSubtree: common.Bool(true),
	}

	// Send the request using the service client
	resp, err := client.ListCategories(context.Background(), req)
	helpers.FatalIfError(err)

	if err != nil {
		return nil, err
	} else {
		var catTmp model.Category
		var listCategory []model.Category

		// Retrieve value from the response.
		for _, s := range resp.Items {
			catTmp = model.Category{*s.Name, *s.Id}
			listCategory = append(listCategory, catTmp)
		}
		return listCategory, nil
	}
}
