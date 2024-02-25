package stravatracker.requestsmodels;

import stravatracker.model.Activity;

import java.util.List;

public class ActivitySumUpResponse {


    private Float totalDistance = 0f;
    private Float averageTime = 0f;
    private Float averageSpeed = 0f;
    private Float averageLength = 0f;
    private Float longestActivity = 0f;
    private Float totalElevationGain = 0f;
    private int totalTrainingsAmount = 0;


    public ActivitySumUpResponse() {
    }

    public ActivitySumUpResponse(List<Activity> activities) {
        this.totalTrainingsAmount = activities.size();

        Float totalTime = 0f;

        for (Activity activity : activities) {
            this.totalDistance += activity.getDistance();
            totalTime += activity.getMovingTime();
            this.totalElevationGain += activity.getTotalElevationGain();
            this.longestActivity = Math.max(this.longestActivity, activity.getDistance());
        }

        this.averageTime = totalTime / activities.size();
        this.averageSpeed = this.totalDistance / totalTime;
        this.averageLength = this.totalDistance / activities.size();

    }


    public Float getTotalDistance() {
        return totalDistance;
    }

    public void setTotalDistance(Float totalDistance) {
        this.totalDistance = totalDistance;
    }

    public Float getAverageTime() {
        return averageTime;
    }

    public void setAverageTime(Float averageTime) {
        this.averageTime = averageTime;
    }

    public Float getAverageSpeed() {
        return averageSpeed;
    }

    public void setAverageSpeed(Float averageSpeed) {
        this.averageSpeed = averageSpeed;
    }

    public Float getAverageLength() {
        return averageLength;
    }

    public void setAverageLength(Float averageLength) {
        this.averageLength = averageLength;
    }

    public Float getLongestActivity() {
        return longestActivity;
    }

    public void setLongestActivity(Float longestActivity) {
        this.longestActivity = longestActivity;
    }

    public Float getTotalElevationGain() {
        return totalElevationGain;
    }

    public void setTotalElevationGain(Float totalElevationGain) {
        this.totalElevationGain = totalElevationGain;
    }

    public int getTotalTrainingsAmount() {
        return totalTrainingsAmount;
    }

    public void setTotalTrainingsAmount(int totalTrainingsAmount) {
        this.totalTrainingsAmount = totalTrainingsAmount;
    }
}
