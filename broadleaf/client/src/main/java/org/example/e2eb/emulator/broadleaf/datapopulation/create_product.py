import requests
import warnings
warnings.filterwarnings('ignore')

header = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36", "Connection": "keep-alive"}


def login(session):
    login_url = "https://localhost:8444/admin/login_admin_post"
    data = { "username": "admin",
             "password": "admin"
            }
    response = session.post(url=login_url, headers=header, data=data, verify=False)

def addProduct(session, num):
    add_url = "https://localhost:8444/admin/product/add"
    data = {
        "ceilingEntityClassname":"org.broadleafcommerce.core.catalog.domain.Product",
        "entityType":"org.broadleafcommerce.core.catalog.domain.ProductImpl",
        "id":"",
        "sectionCrumbs":"",
        "mainEntityName":"",
        "preventSubmit":"false",
        "jsErrorMapString":"",
        "fields['defaultSku__name'].value":"product"+str(num),
        "fields['defaultSku__longDescription'].value":"",
        "fields['defaultCategory'].value":"2001",
        "fields['defaultCategory'].displayValue":"Home",
        "fields['manufacturer'].value":"",
        "fields['url'].value":"/product"+str(num),
        "fields['overrideGeneratedUrl'].value":"false",
        "fields['displayTemplate'].value":"",
        "fields['defaultSku__skuMedia---primary'].value":"",
        "fields['defaultSku__activeStartDate'].value-display":"Monday, April 26, 2021 @ 4:03pm",
        "fields['defaultSku__activeStartDate'].value":"2021.04.26 16:03:52",
        "fields['defaultSku__activeEndDate'].value-display":"",
        "fields['defaultSku__activeEndDate'].value":"",
        "fields['defaultSku__upc'].value":"",
        "fields['defaultSku__externalId'].value":"",
        "fields['metaTitle'].value":"",
        "fields['metaDescription'].value":"",
        "fields['canonicalUrl'].value":"",
        "fields['defaultSku__retailPrice'].value":"10",
        "fields['defaultSku__salePrice'].value":"",
        "fields['defaultSku__cost'].value":"",
        "fields['canSellWithoutOptions'].value":"false",
        "fields['defaultSku__inventoryType'].value":"CHECK_QUANTITY",
        "fields['defaultSku__quantityAvailable'].value":"99999999",
        "fields['defaultSku__dimension__width'].value":"",
        "fields['defaultSku__dimension__height'].value":"",
        "fields['defaultSku__dimension__depth'].value":"",
        "fields['defaultSku__dimension__girth'].value":"",
        "fields['defaultSku__dimension__dimensionUnitOfMeasure'].value":"CENTIMETERS",
        "fields['defaultSku__isMachineSortable'].value":"false",
        "fields['defaultSku__fulfillmentType'].value":"",
        "fields['defaultSku__weight__weight'].value":"",
        "fields['defaultSku__weight__weightUnitOfMeasure'].value":"KILOGRAMS"
    }
    response = session.post(url=add_url, headers=header, data=data, verify=False)

session = requests.Session()
login(session)
for i in range(1,1025):
    addProduct(session, i)
session.close()
