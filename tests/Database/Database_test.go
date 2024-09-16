package Database

import (
	"app/src/Models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositorySuite struct {
	suite.Suite
	testDatabase *TestDatabase
}

func (suite *RepositorySuite) SetupSuite() {
	suite.testDatabase = SetupTestDatabase()
}

func (suite *RepositorySuite) TearDownSuite() {
	suite.testDatabase.TearDown()
}

func (suite *RepositorySuite) TearDownTest() {
	suite.testDatabase.ClearDatabase()
}

func (suite *RepositorySuite) TestInsertActivity() {
	athlete := Models.StravaAthlete{
		Firstname: "John",
		Lastname:  "Doe",
	}
	activity1 := Models.StravaActivity{
		Athlete:            athlete,
		Name:               "Morning Run",
		Distance:           10000,
		MovingTime:         3600,
		TotalElevationGain: 100.0,
		Type:               "Run",
		SportType:          "Running",
	}
	activity2 := Models.StravaActivity{
		Athlete:            athlete,
		Name:               "Evening Run",
		Distance:           15000,
		MovingTime:         5400,
		TotalElevationGain: 150.0,
		Type:               "Run",
	}

	suite.Run("when athlete is in database", func() {
		err := suite.testDatabase.DbService.InsertActivity(activity1)
		suite.Nil(err)
		err = suite.testDatabase.DbService.InsertActivity(activity2)
		suite.Nil(err)
	})

	suite.Run("when athlete is not in database", func() {
		err := suite.testDatabase.DbService.InsertActivity(activity1)
		suite.Nil(err)
		err = suite.testDatabase.DbService.InsertActivity(activity2)
		suite.Nil(err)
	})
}

func (suite *RepositorySuite) TestInsertAthlete() {
	athlete := &Models.StravaAthlete{
		Firstname: "John",
		Lastname:  "Doe",
	}

	suite.Run("when athlete is not in database", func() {
		err := suite.testDatabase.DbService.InsertAthlete(athlete)
		suite.Nil(err)
		suite.NotNil(athlete.ID)
	})

}

func (suite *RepositorySuite) TestGetAthleteIndex() {
	athlete1 := &Models.StravaAthlete{
		Firstname: "John",
		Lastname:  "Doe",
	}

	athlete2 := &Models.StravaAthlete{
		Firstname: "John",
		Lastname:  "Doe",
	}

	suite.Run("when athlete is not in database", func() {
		err := suite.testDatabase.DbService.GetAthleteIndex(athlete1)
		suite.Nil(err)
		suite.NotNil(athlete1.ID)
	})

	suite.Run("when athlete is in database", func() {
		err := suite.testDatabase.DbService.GetAthleteIndex(athlete1)
		suite.Nil(err)
		err = suite.testDatabase.DbService.GetAthleteIndex(athlete2)
		suite.Nil(err)
		suite.Equal(athlete1.ID, athlete2.ID)
	})
}

func (suite *RepositorySuite) TestGetLatestActivity() {
	athlete := Models.StravaAthlete{
		Firstname: "John",
		Lastname:  "Doe",
	}
	activity1 := Models.StravaActivity{
		Athlete:            athlete,
		Name:               "Morning Run",
		Distance:           10000,
		MovingTime:         3600,
		TotalElevationGain: 100.0,
		Type:               "Run",
		SportType:          "Running",
	}
	activity2 := Models.StravaActivity{
		Athlete:            athlete,
		Name:               "Evening Run",
		Distance:           15000,
		MovingTime:         5400,
		TotalElevationGain: 150.0,
		Type:               "Run",
	}

	suite.Run("when there is no activity", func() {
		latestActivity, err := suite.testDatabase.DbService.GetLatestActivity()
		suite.Nil(err)

		suite.Equal(latestActivity.Name, "")
		suite.Equal(latestActivity.Distance, 0.0)
		suite.Equal(latestActivity.MovingTime, 0)
		suite.Equal(latestActivity.TotalElevationGain, 0.0)
		suite.Equal(latestActivity.Type, "")
		suite.Equal(latestActivity.SportType, "")
	})

	suite.Run("when there is activity", func() {

		err := suite.testDatabase.DbService.InsertActivity(activity1)
		suite.Nil(err)
		err = suite.testDatabase.DbService.InsertActivity(activity2)
		suite.Nil(err)
		err = suite.testDatabase.DbService.GetAthleteIndex(&athlete)
		suite.Nil(err)

		latestActivity, err := suite.testDatabase.DbService.GetLatestActivity()

		suite.Nil(err)
		suite.Equal(latestActivity.UserID, athlete.ID)
		suite.Equal(latestActivity.Name, activity2.Name)
		suite.Equal(latestActivity.Distance, activity2.Distance)
		suite.Equal(latestActivity.MovingTime, activity2.MovingTime)
		suite.Equal(latestActivity.TotalElevationGain, activity2.TotalElevationGain)
		suite.Equal(latestActivity.Type, activity2.Type)
		suite.Equal(latestActivity.SportType, activity2.SportType)
	})

}

func (suite *RepositorySuite) TestGetAthletesData() {
	athlete1 := Models.StravaAthlete{
		Firstname: "John",
		Lastname:  "Doe",
	}
	athlete2 := Models.StravaAthlete{
		Firstname: "Jane",
		Lastname:  "Smith",
	}
	activity1 := Models.StravaActivity{
		Athlete:            athlete1,
		Name:               "Morning Run",
		Distance:           10000,
		MovingTime:         3600,
		TotalElevationGain: 100.0,
		Type:               "Run",
		SportType:          "Running",
	}

	activity2 := Models.StravaActivity{
		Athlete:            athlete2,
		Name:               "Evening Run",
		Distance:           15000,
		MovingTime:         5400,
		TotalElevationGain: 150.0,
		Type:               "Run",
		SportType:          "Running",
	}

	activity3 := Models.StravaActivity{
		Athlete:            athlete2,
		Name:               "Evening Run",
		Distance:           15000,
		MovingTime:         5400,
		TotalElevationGain: 150.0,
		Type:               "Run",
		SportType:          "Running",
	}
	suite.Run("when there are no athletes", func() {
		athletesData := suite.testDatabase.DbService.GetAthletesData()
		suite.Len(athletesData, 0)
	})

	suite.Run("when there are athletes with activities", func() {
		err := suite.testDatabase.DbService.InsertAthlete(&athlete1)
		suite.Nil(err)
		err = suite.testDatabase.DbService.InsertAthlete(&athlete2)
		suite.Nil(err)

		err = suite.testDatabase.DbService.InsertActivity(activity1)
		suite.Nil(err)
		err = suite.testDatabase.DbService.InsertActivity(activity2)
		suite.Nil(err)
		err = suite.testDatabase.DbService.InsertActivity(activity3)
		suite.Nil(err)

		athletesData := suite.testDatabase.DbService.GetAthletesData()
		suite.Len(athletesData, 2)

		// Validate athlete1 data
		athleteData1 := athletesData[0]
		suite.Equal(athleteData1.Name, "John Doe")
		suite.Equal(athleteData1.TotalActivities, 1)
		suite.Equal(athleteData1.Distance, activity1.Distance/1000) // in kilometers
		suite.Equal(athleteData1.ElevationGain, activity1.TotalElevationGain)
		suite.Equal(athleteData1.LongestActivity, activity1.Distance/1000) // in kilometers
		suite.Equal(athleteData1.TotalTime, float64(activity1.MovingTime))
		suite.Equal(athleteData1.AverageTime, float64(activity1.MovingTime))
		suite.Equal(athleteData1.AverageLength, activity1.Distance/1000)
		suite.Equal(athleteData1.AverageSpeed, (activity1.Distance/1000)/(float64(activity1.MovingTime)/3600))

		// Validate athlete2 data
		athleteData2 := athletesData[1]
		suite.Equal(athleteData2.Name, "Jane Smith")
		suite.Equal(athleteData2.TotalActivities, 2)
		suite.Equal(athleteData2.Distance, (activity2.Distance+activity3.Distance)/1000) // in kilometers
		suite.Equal(athleteData2.ElevationGain, activity2.TotalElevationGain+activity3.TotalElevationGain)
		suite.Equal(athleteData2.LongestActivity, activity2.Distance/1000) // in kilometers
		suite.Equal(athleteData2.TotalTime, float64(activity2.MovingTime+activity3.MovingTime))
		suite.Equal(athleteData2.AverageTime, float64(activity2.MovingTime+activity3.MovingTime)/2)
		suite.Equal(athleteData2.AverageLength, (activity2.Distance+activity3.Distance)/2000)
		suite.Equal(athleteData2.AverageSpeed, ((activity2.Distance+activity3.Distance)/1000)/(float64(activity2.MovingTime+activity3.MovingTime)/3600))
	})

}

func TestDatabaseModule(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
