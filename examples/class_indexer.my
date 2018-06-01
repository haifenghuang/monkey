//This file demonstrate using of class indexer
//////////////////////////////////////////////////////////////////////////////////
//                                       First
//////////////////////////////////////////////////////////////////////////////////
class IndexedNames
{
    let namelist = []
    let size = 10
    fn init()
    {
        let i = 0
        for (i = 0; i < size; i++)
        {
            namelist[i] = "N. A."
        }
    }

    fn getNameList() {
        println(namelist)
    }

    property this[index]
    {
        get
        {
            let tmp;
            if ( index >= 0 && index <= size - 1 )
            {
               tmp = namelist[index]
            }
            else
            {
               tmp = ""
            }
     
            return tmp
         }
         set
         {
             if ( index >= 0 && index <= size-1 )
             {
                 namelist[index] = value
             }
         }
    }
}

fn Main()
{
    namesObj = new IndexedNames()

    //Calling Indexer's set function
    namesObj[0] = "Zara"
    namesObj[1] = "Riz"
    namesObj[2] = "Nuha"
    namesObj[3] = "Asif"
    namesObj[4] = "Davinder"
    namesObj[5] = "Sunil"
    namesObj[6] = "Rubic"

    namesObj.getNameList()

    for (i = 0; i < namesObj.size; i++)
    {
        println(namesObj[i]) //Calling Indexer's get function
    }
}

Main()

//////////////////////////////////////////////////////////////////////////////////
//                                       Second
//////////////////////////////////////////////////////////////////////////////////
class StudentIdentifications
{
    let studentIDs = [];
    let fullnames = [];

    property this[item]
    {
        get
        {
            if type(item) == "INTEGER" {
                for (i = 0; i < len(studentIDs); i++) {
                    if (item == studentIDs[i]) {
                        return fullnames[i];
                    }
                }
                return "Unknown Student";
            } elseif type(item) == "STRING" {
                for (i = 0; i < len(fullnames); i++) {
                    if (item == fullnames[i]) {
                        return studentIDs[i];
                    }
                }
                return 0;
            }

        }
    }

    fn init()
    {
        
        studentIDs[0] = 39472;
        studentIDs[1] = 13957;
        studentIDs[2] = 73957;
        studentIDs[3] = 97003;
        studentIDs[4] = 28947;
        studentIDs[5] = 97395;

        fullnames[0] = "Paul Bertrand Yamaguchi";
        fullnames[1] = "Ernestine Ngovayang";
        fullnames[2] = "Patricia L Katts";
        fullnames[3] = "Helene Mukoko";
        fullnames[4] = "Joan Ursula Hancock";
        fullnames[5] = "Arlette Mimosa";
    }
}

class Program
{
    static fn Main()
    {
        std = new StudentIdentifications();

        println("Student Identification");
        println("Student ID: 39472");
        printf("Full Name:  %v\n", std[39472]);

        println("Student Identification");
        println("Full Name:  Joan Ursula Hancock");
        printf("Student ID: %v\n", std["Joan Ursula Hancock"]);

        return 0;
    }
}

Program.Main()

//This would produce:
//
//Student Identification
//Student ID: 39472
//Full Name:  Paul Bertrand Yamaguchi
//
//Student Identification
//Full Name:  Joan Ursula Hancock
//Student ID: 28947

//////////////////////////////////////////////////////////////////////////////////
//                Third(Multiple indexer with multiple parameters)
//////////////////////////////////////////////////////////////////////////////////
class CourseScore
{
    let name;
    let courseID;
    let score;

    fn init(name, courseID, score)
    {  
        this.name = name;
        this.courseID = courseID;
        this.score = score;
    }

    property Name
    {
        get { return name; }
        set { name = value; }
    }

    property CourseID
    {
        get { return courseID; }
        set { courseID = value; }
    }

    property Score
    {
        get { return score; }
        set { score = value; }
    }
}

class CourseScoreIndexer
{
    let arrCourseScore = [];

    property this[name, courseID]
    {
        get
        {
            for cs in arrCourseScore
            {
                if (cs.Name == name && cs.CourseID == courseID) {
                    return cs.Score;
                }
            }
            return -1;
        }
        set
        {
            arrCourseScore.push(new CourseScore(name, courseID, value));
        }
    }

    property this[name]
    {
        get
        {
            let tempArray = [];
            for cs in arrCourseScore
            {
                if (cs.Name == name) {
                    tempArray.push(cs);
                }
            }
            return tempArray;
        }
    }
}

class Test
{  
    static fn Main()
    {
        let csi = new CourseScoreIndexer();
        csi["Mike", 1] = 90;
        csi["Mike", 2] = 70;
        csi["Mike", 3] = 85;
        csi["Vicky", 1] = 80;

        println(csi["Mike", 2]); //output : 70
        println("\n<<<<<<<<Mike's all course>>>>>>>>");

        let tempArray = []
        tempArray = csi["Mike"];
        for cs in tempArray
        {
            printf("Name: %v,   CourseID: %v,   Score: %v\n", cs.Name, cs.CourseID, cs.Score);
        }
    }
}

Test.Main()
