package stravatracker.requestsmodels;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDate;

@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
public class DateRange {

    private LocalDate startDate;
    private LocalDate endDate;


    public DateRange(String startDate, String endDate) {
        this.setStartDate(startDate);
        this.setEndDate(endDate);
    }

    public void setStartDate(String startDate) {
        LocalDate parsedStartDate = LocalDate.parse(startDate);
        if (this.endDate != null && parsedStartDate.isAfter(this.endDate)) {
            throw new IllegalArgumentException("Start date cannot be after the end date.");
        }
        this.startDate = parsedStartDate;
    }


    public void setEndDate(String endDate) {
        LocalDate parsedEndDate = LocalDate.parse(endDate);
        if (this.startDate != null && parsedEndDate.isBefore(this.startDate)) {
            throw new IllegalArgumentException("End date cannot be before the start date.");
        }
        this.endDate = parsedEndDate;
    }

    public LocalDate getStartDate() {
        return (this.startDate == null) ? LocalDate.of(1970, 1, 1) : this.startDate;
    }

    public LocalDate getEndDate() {
        return (this.endDate == null) ? LocalDate.now().plusYears(100) : this.endDate;
    }
}


