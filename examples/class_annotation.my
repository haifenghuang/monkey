//This file demonstrate annotation
//////////////////////////////////////////////////////////////////////////////////
//                                       First
//////////////////////////////////////////////////////////////////////////////////
class @Test {
  property Enabled
}

//marker annotation
class @Demo {}


class TestExample {
  @Demo
  @Test{Enabled = true}
  fn TestA() {
    printf("TestA is called\n")
  }

  @Demo
  @Test{Enabled = false}
  fn TestB() {
    printf("TestB is called\n")
  }

  @Demo
  @Test{Enabled = false}
  fn TestC() {
    printf("TestC is called\n")
  }

  @Demo
  @Test{Enabled = true}
  fn TestD() {
    printf("TestD is called\n")
  }
}

testObj = new TestExample()
for method in testObj.getMethods() {
  printf("\nmethodName=%s\n", method.name)
  annos = method.getAnnotations()
  for anno in annos {
    //println()
    //printf("ANNO NAME=%s, enabled=%t\n", anno, anno.Enabled)

    if anno.instanceOf(Test) {
      printf("ANNO NAME=%s, enabled=%t\n", anno, anno.Enabled)
      if anno.Enabled {
        method.invoke()
      }
    } elseif anno.instanceOf(Demo) {
      printf("ANNO NAME=%s, \n", anno)
    }
  }
}



//////////////////////////////////////////////////////////////////////////////////
//                                       Second
//////////////////////////////////////////////////////////////////////////////////
//class @MethodTrace { //Annotation class
//  property LogLevel
//}
//
//enum Logger {
//  INFO,
//  DEBUG
//}
//
//class Calculator {
//  @MethodTrace{LogLevel = Logger.INFO}
//  fn calculatePow(x, y) {
//    result = math.pow(x, y)
//    return result
//  }
//  
//  @MethodTrace{LogLevel = Logger.DEBUG}
//  fn multiply(a, b) {
//    return a * b
//  }
//  
//  @MethodTrace{LogLevel = Logger.DEBUG}
//  fn add(a, b) {
//    return a + b
//  }
//}
//
//class Handler {
//  static fn handle(o) {
//    for m in o.getMethods() {
//      for a in m.getAnnotations() {
//        if (a.is_a(MethodTrace)) {
//          if a.LogLevel == Logger.INFO {
//            println("Performing INFO log for \"" + m.getName() + "\" method: ")
//            printf("calculation result = %v\n\n", m.invoke(2,3))
//          } elseif (a.LogLevel == Logger.DEBUG) {
//            println("Performing DEBUG log for \"" + m.getName() + "\" method")
//            printf("calculation result = %d\n\n", m.invoke(10,5))
//          }
//        }
//      }
//    }
//  }
//}
//
//class Main {
//  static fn main() {
//    Calculator calculator = new Calculator()
//    Handler.handle(calculator)
//  }
//}
//
//Main.main()


//////////////////////////////////////////////////////////////////////////////////
//                                       Third
//////////////////////////////////////////////////////////////////////////////////
//Annotations
//class @Test {}
//class @TestSetup {}
//class @TestTearDown {}
//
////Calculator class
//class Calculator{
//  @TestSetup
//  static fn Setup() {
//    println("TestSetup called!\n")
//  }
//
//  @TestTearDown
//  static fn TearDown() {
//    println("\nTestTearDown called!\n")
//  }
//
//  @Test
//  fn add(x, y) {
//    return x + y
//  }
//
//  @Test
//  fn sub(x, y) {
//    return x - y
//  }
//
//  @Test
//  fn mul(x, y) {
//    return x * y
//  }
//
//  @Test
//  fn div(x, y) {
//    return x / y
//  }
//}
//
//class Handler {
//  static fn handle(o) {
//    methods = o.getMethods()
//
//    //run test setup
//    for m in methods {
//      testSetupAnno = m.getAnnotation(TestSetup)
//      if testSetupAnno != nil {
//        m.invoke()
//        break
//      }
//    }
//
//    //run test
//    for m in methods {
//      testAnno = m.getAnnotation(Test)
//      if testAnno != nil {
//        printf("%s(20, 10) = %v\n", m.name, m.invoke(20, 10))
//      }
//    }
//
//    //run test teardown
//    for m in methods {
//      testTearDownAnno = m.getAnnotation(TestTearDown)
//      if testTearDownAnno != nil {
//        m.invoke()
//        break
//      }
//    }
//  }
//}
//
//class Main {
//  static fn main() {
//    Calculator calculator = new Calculator()
//    Handler.handle(calculator)
//  }
//}
//
//Main.main()

////////////////////////////////////////////////////////////////////////////////
//                                       Fourth
////////////////////////////////////////////////////////////////////////////////
//class @ParentMinMaxValidator {
//  property MaxLength default 10
//}
//
//class @MinMaxValidator : ParentMinMaxValidator {
//  property MinLength
//  //property MaxLength default 10
//}
//
//class @NoSpaceValidator {}
//
//class @DepartmentValidator {
//  property Department
//}
//
//class Request {
//  //@MinMaxValidator(MinLength=1, MaxLength=10)
//  @MinMaxValidator(MinLength=1)
//  property FirstName;
//
//  @NoSpaceValidator
//  property LastName;
//
//  @DepartmentValidator(Department=["Department of Education", "Department of Labors", "Department of Justice"])
//  property Dept;
//}
//
//class RequestHandler {
//  static fn handle(o) {
//    props = o.getProperties()
//    for p in props {
//      annos = p.getAnnotations()
//      for anno in annos {
//        if anno.instanceOf(MinMaxValidator) {
//          if len(p.value) > anno.MaxLength || len(p.value) < anno.MinLength {
//            printf("Property '%s' is not valid!\n", p.name)
//          }
//        } elseif anno.instanceOf(NoSpaceValidator) {
//          for c in p.value {
//            if c == " " || c == "\t" {
//              printf("Property '%s' is not valid!\n", p.name)
//              break
//            }
//          }
//        } elseif anno.instanceOf(DepartmentValidator) {
//          found = false
//          for d in anno.Department {
//            if p.value == d {
//              found = true
//            }
//          }
//          if !found {
//            printf("Property '%s' is not valid!\n", p.name)
//          }
//        }
//      }
//    }
//  }
//}
//
//class RequestMain {
//  static fn main() {
//    request = new Request()
//    request.FirstName = "Haifeng123456789" //not valid
//    request.LastName = "Huang    "  //not valid
//    request.Dept = "Department of Justice"
//    RequestHandler.handle(request)
//  }
//}
//
//RequestMain.main()
