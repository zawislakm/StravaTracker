package stravatracker.requestsmodels;

import stravatracker.model.SportType;

public class ActivitySumUpRequest {
    private DateRange dateRange;
    private SportType sportType;


    public ActivitySumUpRequest() {
    }


    public ActivitySumUpRequest(DateRange dateRange) {
        this.dateRange = dateRange;
    }

    public ActivitySumUpRequest(SportType sportType) {
        this.sportType = sportType;
    }

    public ActivitySumUpRequest(DateRange dateRange, SportType sportType) {
        this.dateRange = dateRange;
        this.sportType = sportType;
    }

    public DateRange getDateRange() {
        return dateRange;
    }

    public void setDateRange(DateRange dateRange) {
        this.dateRange = dateRange;
    }

    public SportType getSportType() {
        return sportType;
    }

    public void setSportType(SportType sportType) {
        this.sportType = sportType;
    }

    @Override
    public String toString() {
        return "ActivitySumUpRequest{" +
                "dateRange=" + dateRange +
                ", sportType=" + sportType +
                '}';
    }
}
