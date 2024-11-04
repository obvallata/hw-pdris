from src.lev_space import levenshtein


def get_close_animals(name):
    dogs = open("/app/src/dogs.txt", 'r')
    cats = open("/app/src/cats.txt", 'r')
    min_space_dog = len(name) + 1
    close_dog = 'not found'
    dog = 'dog'
    while len(dog) != 0:
        dog = dogs.readline()
        new_space = levenshtein(name.lower().replace('\n', ''), dog.lower().replace('\n', ''))
        if new_space < min_space_dog:
            min_space_dog = new_space
            close_dog = dog[::]
    min_space_cat = len(name) + 1
    close_cat = 'not found'
    cat = 'cat'
    while len(cat) != 0:
        cat = cats.readline()
        new_space = levenshtein(name.lower().replace('\n', ''), cat.lower().replace('\n', ''))
        if new_space < min_space_cat:
            min_space_cat = new_space
            close_cat = cat[::]
    return close_dog.replace('\n', ''), min_space_dog, close_cat.replace('\n', ''), min_space_cat
