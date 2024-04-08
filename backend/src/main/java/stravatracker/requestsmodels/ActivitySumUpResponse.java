package stravatracker.requestsmodels;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import stravatracker.model.Activity;

import java.util.List;

@NoArgsConstructor
@AllArgsConstructor
@Setter
@Getter
public class ActivitySumUpResponse {


    private Float totalDistance = 0f;
    private Float averageTime = 0f;
    private Float averageSpeed = 0f;
    private Float averageLength = 0f;
    private Float longestActivity = 0f;
    private Float totalElevationGain = 0f;
    private int totalTrainingsAmount = 0;


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


}
