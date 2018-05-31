//This file demonstrate some common class operation
//////////////////////////////////////////////////////////////////////////////////
//                                       First
//////////////////////////////////////////////////////////////////////////////////
//let Animal = class {
//	let name = "Animal"
//
//	fn init(msg) {
//		name = msg
//		printf("%s\n", msg)
//	}
//
//	fn Walk1(msg) {
//		printf("Animal Walk:%s\n", msg)
//	}
//}
//
//class Dog : Animal {
//	let name = "NONAME"
//
//	fn init(name) {
//		parent.init("Calling parent")
//		this.name = name
//	}
//
//	fn Walk() {
//		printf("Hello, my name is %s\n", name)
//	}
//}
//
//let dogObj = new Dog("doggy")
//dogObj.Walk()
//
//printf("Dog.name=%s\n", dogObj.name)
//dogObj.name = "honey"
//printf("Dog.name=%s\n", dogObj.name)
//
//dogObj.Walk1("dog")



//////////////////////////////////////////////////////////////////////////////////
//                                       Second
//////////////////////////////////////////////////////////////////////////////////
//class Animal {
//    let Name;
//
//    fn MakeNoise()
//    {
//        println("generic noise")
//    }
//    fn ToString()
//    {
//        return "oooooooo"
//    }
//}
//
//class Cat : Animal {
//    fn init(name)
//    {
//        this.Name = name
//    }
//
//    fn MakeNoise()
//    {
//        println("Meow")
//    }
//
//    fn ToString()
//    {
//        return Name + " cat"
//    }
//}
//
//class Dog : Animal {
//    fn init(name)
//    {
//        this.Name = name
//    }
//
//    fn MakeNoise()
//    {
//        println("Woof!")
//    }
//
//    fn ToString()
//    {
//        return Name + " dog"
//    }
//
//    fn OnlyDogMethod()
//    {
//        println("secret dog only method")
//    }
//}
//
//
//cat = new Cat("pearl")
//dog = new Dog("cole")
//randomAnimal = new Animal()
//
//println()
//animals = [cat, dog, randomAnimal]
//for animal in animals
//{
//    println("Animal name: " + animal.Name)
//    animal.MakeNoise()
//    println(animal.ToString())
//
//    if animal.is_a(Dog) {
//        animal.OnlyDogMethod()
//    }
//    println()
//    println()
//}
//
//
////let testAnimal = fn(animal) {
////    println("Animal name: " + animal.Name)
////    animal.MakeNoise()
////    println(animal.ToString())
////
////    if animal.is_a(Dog) {
////        animal.OnlyDogMethod()
////    }
////    println()
////    println()
////}
////
////testAnimal(cat)
////testAnimal(dog)
////testAnimal(randomAnimal)

//////////////////////////////////////////////////////////////////////////////////
//                                       Third
//////////////////////////////////////////////////////////////////////////////////
class Animal
{
    let m_name;
    let m_voice;
    
    fn init(name)
    {
        this.m_name = name;
        this.m_voice = "???";
    }
    

    fn getName() { return m_name; }
    fn setName(name) { m_name = name; }
    fn getVoice() { return m_voice; }
    fn setVoice(voice) { m_voice = voice; }

    fn speak() 
    {
        println(m_name + " has voice as " + m_voice);
        println(m_name + " speaks '???! ???! ???!'");
    }

    fn eat()
    {
        println(m_name + " likes eating ???");
    }
}


class Cat : Animal {
    fn init(name)
    {
        parent.init(name)
        m_voice = "Meo";
    }

    fn speak()  
    {
        println(m_name + " has voice as " + m_voice);
        println(m_name + " speaks 'Meo! Meo! Meo!'");
    }

    fn eat()
    {
        println(m_name + " likes eating Mouse");
    }
}

class Dog : Animal {

    fn init(name)
    {
        parent.init(name)
        m_voice = "Gau";
    }

    fn speak()  
    {
        println(m_name + " has voice as " + m_voice);
        println(m_name + " speaks 'Gau! Gau! Gau!'");
    }

    fn eat()
    {
        println(m_name + " likes eating Cat");
    }
}

myAnimal = new Animal("My Animal");
myAnimal.speak();

myDog = new Dog("My Dog");
myDog.speak();

myCat = new Cat("My Cat");
myCat.speak();

myAnimalDog = new Dog("My Annimal Dog");
myAnimalDog.speak();

myAnimalCat = new Cat("My Animal Cat");
myCat.eat();
myDog.eat();

myAnimalDog.eat();
myAnimalCat.eat();

println(myAnimalCat.toString())        //string() method is coming from the root class 'object'
println(myAnimalCat.hashCode())        //hashCode() method is coming from the root class 'object'
println(myAnimalCat.is_a(Cat))         //is_a() method is coming from the root class 'object'
println(myAnimalCat.instanceOf(Cat))   //instanceOf() method is is equal to is_a() method
println(myAnimalCat.classOf())         //classOf() method is coming from the root class 'object'


//////////////////////////////////////////////////////////////////////////////////
//                                       Fourth
//////////////////////////////////////////////////////////////////////////////////
//class Person {
//	let firstName;
//	let lastName;
//	let nickName;
//
//	fn init(firstname, lastname="huang", nickname="mike") {
//		this.firstName = firstname
//		this.lastName = lastname
//		this.nickName = nickname
//	}
//
//	fn Msg(extra1, extra2="extra2", args...) {
//		printf("firstName=%s, lastName=%s, nickName=%s\n", firstName, lastName, nickName)
//		printf("extra1=%s, extra2=%s\n", extra1, extra2)
//		for i in args {
//			printf("i=%v\n", i)
//		}
//	}
//}
//
//obj = new Person("haifeng")
//obj.Msg("extra1", "Hello", 1, [1,2,3], (4,5,6))
//
//println()
//obj.Msg("extra1")
