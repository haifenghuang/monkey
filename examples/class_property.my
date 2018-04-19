//This file demonstrate class property
class Date {
    let month = 7;  // Backing store
    property Month
    {
        get { return month }
        set {
            if ((value > 0) && (value < 13))
            {
                month = value
            } else {
               println("BAD, month is invalid")
            }
        }
    }

    property Year;
    //property Year { get; set;}

    property Day { get; }

    property OtherInfo1 { get; }
    property OtherInfo2 { set; }

    fn init(year, month, day) {
        this.Year = year
        this.Month = month
        this.Day = day
    }

    fn getDateInfo() {
        printf("Year:%v, Month:%v, Day:%v\n", this.Year, this.Month, this.Day) //note here, you need to use 'this.Property', not 'Property'
    }


}

dateObj = new Date(2000, 5, 11)
//printf("Calling Date's getter, month=%d\n", dateObj.Month)
dateObj.getDateInfo()

println()
dateObj.Month = 10
printf("dateObj.Month=%d\n", dateObj.Month)

dateObj.Year = 2018
println()
dateObj.getDateInfo()

//Below code will raise an execution error! Because OtherInfo1 is a READONLY property.
//dateObj.OtherInfo1 = "Other Date Info"
//println(dateObj.OtherInfo1)

//Below code will raise an execution error! Because OtherInfo2 is a WRITEONLY property.
//dateObj.OtherInfo2 = "Other Date Info2"
//println(dateObj.OtherInfo2)

//Below code will raise an execution error! Because Day is a READONLY property.
//dateObj.Day = 18
