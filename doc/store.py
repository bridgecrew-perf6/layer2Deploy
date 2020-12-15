OntCversion = '2.0.0'
from ontology.interop.System.Storage import Put, GetContext, Get
from ontology.interop.System.Runtime import Notify,CheckWitness
from ontology.interop.Ontology.Contract import Migrate
from ontology.interop.System.Contract import Destroy

KeyOwnerAddress = "OwnerAddress"

def Main(operation, args):
    if operation == "init":
        return init(args)
    elif operation == "init_status":
        return init_status()
    elif operation == 'StoreHash':
        return StoreHash(args)
    elif operation == "CheckHashExist":
        return CheckHashExist(args)
    elif operation == "destroyContract":
        return DestroyContract()
    elif operation == "migrateContract":
        return MigrateContract(args)
    return False


def DestroyContract():
    addr = Get(GetContext(), KeyOwnerAddress)
    assert(len(addr) != 0)
    assert(CheckWitness(addr))
    return Destroy()


def MigrateContract(code):
    addr = Get(GetContext(), KeyOwnerAddress)
    assert(len(addr) != 0)
    assert(CheckWitness(addr))

    success = Migrate(code, True, "name", "version", "author", "email", "description")
    assert(success)
    Notify(["Migrate successfully", success])
    return success


def init_status():
    addr = Get(GetContext(), KeyOwnerAddress)
    Notify([addr])
    return addr


def init(addr):
    if len(Get(GetContext(), KeyOwnerAddress)) == 0:
        Put(GetContext(), KeyOwnerAddress, addr)
        Notify(["init True"])
        return True
    else:
        Notify(["init False"])
        return False


def StoreHash(args):
    if len(args) != 1:
        return False

    addr = Get(GetContext(), KeyOwnerAddress)
    assert(len(addr) != 0)
    assert(CheckWitness(addr))

    inputHash = args[0]

    Put(GetContext(), inputHash, 1)
    Notify([inputHash])
    return True


def CheckHashExist(args):
    if len(args) != 0:
        return False

    inputHash = args[0]
    exist = Get(GetContext(), inputHash)
    return exist + 0
