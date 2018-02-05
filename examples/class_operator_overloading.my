//This file demonstrate operator overloading
class Vector {
	// instance variables
	let x = 0;
	let y = 0;
	let z = 0;

	// constructor
	fn init (a, b, c) {
		if (!a) { a = 0;}
		if (!b) {b = 0;}
		if (!c) {c = 0;}
		x = a; y = b; z = c;
	}

	// instance method (built-in operator overriding)
	fn +(v) {
		if (type(v) == "INTEGER") {
			return new Vector(x+v, y+v, z+v);
		}
		elseif v.is_a(Vector) {
			return new Vector(x+v.x, y+v.y, z+v.z);
		}
		return nil;
	}

	fn -() {
		return new Vector(-(x),-(y),-(z))
	}

	fn ++() {
		return new Vector(x+1,y+1,z+1)
	}

	fn --() {
		return new Vector(x-1,y-1,z-1)
	}

	// instance method (built-in String conversion overriding)
	fn String() {
			// string interpolation support
		return fmt.sprintf("(%v),(%v),(%v)", this.x, this.y, this.z);
	}

	fn otherFunc(v) {
		return this + v  //calling the '+' function
	}
}

fn Vectormain() {
	// initialize a new vector object
	v1 = new Vector(1,2,3);
	// initialize a new vector object
	v2 = new Vector(4,5,6);

	// call + function in the vector object
	v3 = v1 + v2 //same as 'v3 = v1.+(v2)'
	// returns string "(5),(7),(9)"
	println(v3.String());

	v4 =v1 + 10 //same as v4 = v1.+(10);
	//returns string "(11),(12),(13)"
	println(v4.String());
	println(v1.otherFunc(8).String())

	v5 = -v4
	println(v5.String())

	v6 = --v4
	println(v6.String())

	v7 = ++v4
	println(v7.String())
}

Vectormain()


class Box {
	let length;
	let breadth;
	let height;

	fn getVolume() {
		return length * breadth * height;
	}

	fn setLength(len) {
		length = len;
	}

	fn setBreadth(bre) {
		breadth = bre;
	}
 
	fn setHeight(hei) {
		height = hei;
	}
	
	// Overload + operator to add two Box objects.
	fn +(other) {
		box = new Box();
		box.length = this.length + other.length;
		box.breadth = this.breadth + other.breadth;
		box.height = this.height + other.height;
		return box;
	}
}

fn Boxmain() {
	// box 1 specification
	Box1 = new Box();
	Box1.setLength(6.0);
	Box1.setBreadth(7.0);
	Box1.setHeight(5.0);

	// box 2 specification
	Box2 = new Box();
	Box2.setLength(12.0);
	Box2.setBreadth(13.0);
	Box2.setHeight(10.0);

	// volume of box 1
	volume = Box1.getVolume();
	printf("Volume of Box1 : %v\n", volume);

	// volume of box 2
	volume = Box2.getVolume();
	printf("Volume of Box2 : %v\n", volume);

	// Add two object as follows:
	Box3 = Box1 + Box2;

	// volume of box 3
	volume = Box3.getVolume();
	printf("Volume of Box3 : %v\n", volume);
}

Boxmain()