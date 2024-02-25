package stravatracker.requestsmodels;

import java.time.LocalDate;

public class DateRange {

    private LocalDate startDate;
    private LocalDate endDate;

    public DateRange() {
    }

    public DateRange(String startDate, String endDate) {
        this.setStartDate(startDate);
        this.setEndDate(endDate);
    }

    public DateRange(LocalDate startDate, LocalDate endDate) {
        this.startDate = startDate;
        this.endDate = endDate;
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
        if (this.startDate == null) {
            return LocalDate.of(1970, 1, 1);
        }
        return this.startDate;
    }

    public LocalDate getEndDate() {
        if (this.endDate == null) {
            return LocalDate.now().plusYears(100);
        }
        return this.endDate;
    }
}


