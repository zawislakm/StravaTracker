package stravatracker.requestsmodels;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import stravatracker.model.SportType;


@NoArgsConstructor
@AllArgsConstructor
@Setter
@Getter
public class ActivitySumUpRequest {
    private DateRange dateRange;
    private SportType sportType;


    public ActivitySumUpRequest(DateRange dateRange) {
        this.dateRange = dateRange;
    }

    public ActivitySumUpRequest(SportType sportType) {
        this.sportType = sportType;
    }


}
