package service

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"git.codenrock.com/cnrprod1725727333-user-88349/zadanie-6105/internal/tender/domain"
// 	"git.codenrock.com/cnrprod1725727333-user-88349/zadanie-6105/internal/tender/domain/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreateTender(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	err:= errors.New("")

// 	m := mocks.NewMockTenderRepositoryGetter(ctrl)

// 	m.EXPECT().CreateTender(gomock.Any(), gomock.Any()).Return(domain.Tender{},err).Times(1)
// 	m.EXPECT().CreateTender(gomock.Any(), gomock.Any()).Return(domain.Tender{},nil).Times(1)

// 	tenderService := NewTenderGetter(m)
// 	testCases:=  []struct{
// 		name string
// 		tender domain.Tender
// 		wantErr error
// 	} {{
// 		name: "error from repository",
// 		tender: domain.Tender{
// 			ID :"string",    
// 			Name: "string" ,   
// 			Description:    " string",
// 			Status  :        " string ",
// 			ServiceType:     " string ", 
// 			OrganizationId:  " string ",   
// 			CreatorUsername: " string ",  
// 			Version :        5,      
// 			CreatedAt:       time.Now(),
// 		},
// 		wantErr: err,
// 	},
// 	{
// 		name: "success",
// 		tender: domain.Tender{
// 			ID :"string",    
// 			Name: "string" ,   
// 			Description:    " string",
// 			Status  :        " string ",
// 			ServiceType:     " string ", 
// 			OrganizationId:  " string ",   
// 			CreatorUsername: " string ",  
// 			Version :        5,      
// 			CreatedAt:       time.Now(),
// 		},
// 		wantErr: nil,

// 	},
// }

// for _, testCase :=range testCases{
// 	t.Run(testCase.name, func(t *testing.T){
// 		_,errTender:=tenderService.CreateTender(context.Background(), testCase.tender)

// 		if testCase.wantErr!= nil{
// 			require.EqualError(t,errTender,"service.CreateTender: "+testCase.wantErr.Error())
// 		} else{
// 			require.NoError(t,errTender)
// 		}

// 	})


// }

// }


